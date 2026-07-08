package service

import (
	"context"
	"errors"
	"html"
	"math"
	"time"

	"smartfarming/dto"
	"smartfarming/model"
	"smartfarming/repository"

	"github.com/google/uuid"
)

type CategoryService interface {
	Create(ctx context.Context, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.CategoryResponse, error)
	GetBySlug(ctx context.Context, slug string) (*dto.CategoryResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, req dto.PaginationRequest) (*dto.PaginationResponse, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) Create(ctx context.Context, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	slug := Slugify(req.Name)
	existing, _ := s.categoryRepo.GetBySlug(ctx, slug)
	if existing != nil {
		return nil, errors.New("category name or slug already exists")
	}

	category := &model.Category{
		Name:        html.EscapeString(req.Name),
		Slug:        slug,
		Description: html.EscapeString(req.Description),
	}

	err := s.categoryRepo.Create(ctx, category)
	if err != nil {
		return nil, err
	}

	res := mapToCategoryResponse(category)
	return &res, nil
}

func (s *categoryService) GetByID(ctx context.Context, id uuid.UUID) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := mapToCategoryResponse(category)
	return &res, nil
}

func (s *categoryService) GetBySlug(ctx context.Context, slug string) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	res := mapToCategoryResponse(category)
	return &res, nil
}

func (s *categoryService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		escapedName := html.EscapeString(*req.Name)
		category.Name = escapedName
		slug := Slugify(escapedName)
		if slug != category.Slug {
			existing, _ := s.categoryRepo.GetBySlug(ctx, slug)
			if existing != nil {
				return nil, errors.New("category name or slug already exists")
			}
			category.Slug = slug
		}
	}

	if req.Description != nil {
		category.Description = html.EscapeString(*req.Description)
	}

	err = s.categoryRepo.Update(ctx, category)
	if err != nil {
		return nil, err
	}

	res := mapToCategoryResponse(category)
	return &res, nil
}

func (s *categoryService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.categoryRepo.Delete(ctx, id)
}

func (s *categoryService) List(ctx context.Context, req dto.PaginationRequest) (*dto.PaginationResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	categories, total, err := s.categoryRepo.List(ctx, req.Page, req.Limit, req.Search)
	if err != nil {
		return nil, err
	}

	var data []dto.CategoryResponse
	for _, c := range categories {
		data = append(data, mapToCategoryResponse(&c))
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

func mapToCategoryResponse(category *model.Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		CreatedAt:   category.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   category.UpdatedAt.Format(time.RFC3339),
	}
}
