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

type ArticleService interface {
	Create(ctx context.Context, req dto.CreateArticleRequest, userID uuid.UUID) (*dto.ArticleResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.ArticleResponse, error)
	GetBySlug(ctx context.Context, slug string) (*dto.ArticleResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateArticleRequest) (*dto.ArticleResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, req dto.PaginationRequest, categoryID *uuid.UUID) (*dto.PaginationResponse, error)
	UpdateImage(ctx context.Context, id uuid.UUID, imageURL string) (*dto.ArticleResponse, error)
}

type articleService struct {
	articleRepo  repository.ArticleRepository
	categoryRepo repository.CategoryRepository
}

func NewArticleService(articleRepo repository.ArticleRepository, categoryRepo repository.CategoryRepository) ArticleService {
	return &articleService{
		articleRepo:  articleRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *articleService) Create(ctx context.Context, req dto.CreateArticleRequest, userID uuid.UUID) (*dto.ArticleResponse, error) {
	// Verify category exists
	category, err := s.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		return nil, errors.New("invalid category ID: category not found")
	}

	slug := Slugify(req.Title)
	// Handle potential duplicate slug by appending unique string if needed
	existing, _ := s.articleRepo.GetBySlug(ctx, slug)
	if existing != nil {
		slug = slug + "-" + uuid.New().String()[:8]
	}

	article := &model.Article{
		Title:      html.EscapeString(req.Title),
		Slug:       slug,
		Content:    html.EscapeString(req.Content),
		CategoryID: req.CategoryID,
	}
	article.CreatedBy = &userID

	err = s.articleRepo.Create(ctx, article)
	if err != nil {
		return nil, err
	}

	// Fetch newly created article to ensure User and Category relationships are preloaded correctly
	resArticle, err := s.articleRepo.GetByID(ctx, article.ID)
	if err != nil {
		// Fallback map if preload fetch fails (e.g. in simple mocks)
		article.Category = *category
		res := mapToArticleResponse(article)
		return &res, nil
	}

	res := mapToArticleResponse(resArticle)
	return &res, nil
}

func (s *articleService) GetByID(ctx context.Context, id uuid.UUID) (*dto.ArticleResponse, error) {
	article, err := s.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := mapToArticleResponse(article)
	return &res, nil
}

func (s *articleService) GetBySlug(ctx context.Context, slug string) (*dto.ArticleResponse, error) {
	article, err := s.articleRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	res := mapToArticleResponse(article)
	return &res, nil
}

func (s *articleService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateArticleRequest) (*dto.ArticleResponse, error) {
	article, err := s.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		escapedTitle := html.EscapeString(*req.Title)
		article.Title = escapedTitle
		slug := Slugify(escapedTitle)
		if slug != article.Slug {
			existing, _ := s.articleRepo.GetBySlug(ctx, slug)
			if existing != nil {
				slug = slug + "-" + uuid.New().String()[:8]
			}
			article.Slug = slug
		}
	}

	if req.Content != nil {
		article.Content = html.EscapeString(*req.Content)
	}

	if req.CategoryID != nil {
		category, err := s.categoryRepo.GetByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, errors.New("invalid category ID: category not found")
		}
		article.CategoryID = *req.CategoryID
		article.Category = *category
	}

	err = s.articleRepo.Update(ctx, article)
	if err != nil {
		return nil, err
	}

	// Refetch to ensure preloads are fresh
	resArticle, err := s.articleRepo.GetByID(ctx, article.ID)
	if err != nil {
		res := mapToArticleResponse(article)
		return &res, nil
	}

	res := mapToArticleResponse(resArticle)
	return &res, nil
}

func (s *articleService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.articleRepo.Delete(ctx, id)
}

func (s *articleService) List(ctx context.Context, req dto.PaginationRequest, categoryID *uuid.UUID) (*dto.PaginationResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	articles, total, err := s.articleRepo.List(ctx, req.Page, req.Limit, req.Search, categoryID)
	if err != nil {
		return nil, err
	}

	var data []dto.ArticleResponse
	for _, a := range articles {
		data = append(data, mapToArticleResponse(&a))
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

func mapToArticleResponse(article *model.Article) dto.ArticleResponse {
	return dto.ArticleResponse{
		ID:        article.ID,
		Title:     article.Title,
		Slug:      article.Slug,
		Content:   article.Content,
		Category:  mapToCategoryResponse(&article.Category),
		ImageURL:  article.ImageURL,
		CreatedBy: article.CreatedBy,
		CreatedAt: article.CreatedAt.Format(time.RFC3339),
		UpdatedAt: article.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *articleService) UpdateImage(ctx context.Context, id uuid.UUID, imageURL string) (*dto.ArticleResponse, error) {
	article, err := s.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	article.ImageURL = imageURL
	err = s.articleRepo.Update(ctx, article)
	if err != nil {
		return nil, err
	}

	res := mapToArticleResponse(article)
	return &res, nil
}
