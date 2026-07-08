package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"smartfarming/config"
	"smartfarming/dto"
	"smartfarming/model"
	"smartfarming/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	GetMe(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
	Logout(ctx context.Context, tokenStr string) error
	VerifyOTP(ctx context.Context, req dto.VerifyOTPRequest) (*dto.AuthResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
}

type tempRegisterData struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	OTP      string `json:"otp"`
}

var (
	tempRegisterStore = make(map[string]tempRegisterData)
	tempStoreMutex   sync.RWMutex
)

func generateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email is already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	otp, err := generateOTP()
	if err != nil {
		return nil, err
	}

	tempData := tempRegisterData{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		OTP:      otp,
	}

	log.Printf("Generated OTP for %s: %s", req.Email, otp)

	if config.RedisClient != nil {
		jsonData, err := json.Marshal(tempData)
		if err != nil {
			return nil, err
		}
		redisKey := "otp:register:" + req.Email
		err = config.RedisClient.Set(ctx, redisKey, jsonData, 5*time.Minute).Err()
		if err != nil {
			return nil, errors.New("failed to save registration OTP session")
		}
	} else {
		tempStoreMutex.Lock()
		tempRegisterStore[req.Email] = tempData
		tempStoreMutex.Unlock()
	}

	return &dto.RegisterResponse{
		Message: "Registration initiated. Verification OTP code has been generated.",
		OTP:     otp,
	}, nil
}

func (s *authService) VerifyOTP(ctx context.Context, req dto.VerifyOTPRequest) (*dto.AuthResponse, error) {
	var tempData tempRegisterData
	found := false

	if config.RedisClient != nil {
		redisKey := "otp:register:" + req.Email
		val, err := config.RedisClient.Get(ctx, redisKey).Result()
		if err == nil {
			err = json.Unmarshal([]byte(val), &tempData)
			if err == nil {
				found = true
			}
		}
	} else {
		tempStoreMutex.RLock()
		data, exists := tempRegisterStore[req.Email]
		if exists {
			tempData = data
			found = true
		}
		tempStoreMutex.RUnlock()
	}

	if !found {
		return nil, errors.New("registration session not found or expired")
	}

	if tempData.OTP != req.Code {
		return nil, errors.New("invalid verification OTP code")
	}

	user := &model.User{
		Name:     tempData.Name,
		Email:    tempData.Email,
		Password: tempData.Password,
		Role:     "user",
	}

	err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	if config.RedisClient != nil {
		redisKey := "otp:register:" + req.Email
		config.RedisClient.Del(ctx, redisKey)
	} else {
		tempStoreMutex.Lock()
		delete(tempRegisterStore, req.Email)
		tempStoreMutex.Unlock()
	}

	token, err := s.generateJWT(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User:  mapToUserResponse(user),
	}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := s.generateJWT(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User:  mapToUserResponse(user),
	}, nil
}

func (s *authService) GetMe(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := mapToUserResponse(user)
	return &res, nil
}

func (s *authService) generateJWT(userID uuid.UUID, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID.String(),
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // Expires in 24h
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

func (s *authService) Logout(ctx context.Context, tokenStr string) error {
	claims := jwt.MapClaims{}
	_, _, err := jwt.NewParser().ParseUnverified(tokenStr, claims)
	if err != nil {
		return errors.New("invalid token format")
	}

	expVal, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("invalid token expiration claim")
	}

	expTime := time.Unix(int64(expVal), 0)
	remainingTime := time.Until(expTime)

	if remainingTime <= 0 {
		return nil
	}

	if config.RedisClient != nil {
		blacklistKey := "blacklist:" + tokenStr
		err = config.RedisClient.Set(ctx, blacklistKey, "1", remainingTime).Err()
		if err != nil {
			return errors.New("failed to blacklist token")
		}
	}

	return nil
}

func mapToUserResponse(user *model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		PhotoURL:  user.PhotoURL,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}
