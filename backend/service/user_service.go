package service

import (
	"context"
	"errors"
	"math"

	"smartfarming/dto"
	"smartfarming/model"
	"smartfarming/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, req dto.PaginationRequest) (*dto.PaginationResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
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
		Role:     req.Role,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	res := mapToUserResponse(user)
	return &res, nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	res := mapToUserResponse(user)
	return &res, nil
}

func (s *userService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		if *req.Email != user.Email {
			existing, _ := s.userRepo.GetByEmail(ctx, *req.Email)
			if existing != nil {
				return nil, errors.New("email is already in use by another user")
			}
			user.Email = *req.Email
		}
	}
	if req.Password != nil && *req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashed)
	}
	if req.Role != nil {
		user.Role = *req.Role
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	res := mapToUserResponse(user)
	return &res, nil
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) List(ctx context.Context, req dto.PaginationRequest) (*dto.PaginationResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	users, total, err := s.userRepo.List(ctx, req.Page, req.Limit, req.Search)
	if err != nil {
		return nil, err
	}

	var data []dto.UserResponse
	for _, u := range users {
		data = append(data, mapToUserResponse(&u))
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	return &dto.PaginationResponse{
		Data: data,
		Meta: dto.Meta{
			CurrentPage:  req.Page,
			Limit:        req.Limit,
			TotalRecords: total,
			TotalPages:   totalPages,
		},
	}, nil
}
