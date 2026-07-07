package service

import (
	"context"
	"errors"
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
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	GetMe(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
	Logout(ctx context.Context, tokenStr string) error
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email is already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
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
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}
