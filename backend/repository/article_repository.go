package repository

import (
	"context"
	"smartfarming/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(ctx context.Context, article *model.Article) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Article, error)
	GetBySlug(ctx context.Context, slug string) (*model.Article, error)
	Update(ctx context.Context, article *model.Article) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page int, limit int, search string, categoryID *uuid.UUID) ([]model.Article, int64, error)
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) Create(ctx context.Context, article *model.Article) error {
	return r.db.WithContext(ctx).Create(article).Error
}

func (r *articleRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Article, error) {
	var article model.Article
	err := r.db.WithContext(ctx).
		Preload("Category").
		First(&article, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) GetBySlug(ctx context.Context, slug string) (*model.Article, error) {
	var article model.Article
	err := r.db.WithContext(ctx).
		Preload("Category").
		First(&article, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) Update(ctx context.Context, article *model.Article) error {
	return r.db.WithContext(ctx).Save(article).Error
}

func (r *articleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	article, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Delete(article).Error
}

func (r *articleRepository) List(ctx context.Context, page int, limit int, search string, categoryID *uuid.UUID) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Article{})

	if categoryID != nil && *categoryID != uuid.Nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("title ILIKE ? OR content ILIKE ?", searchPattern, searchPattern)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Offset(offset).
		Limit(limit).
		Order("created_at desc").
		Preload("Category").
		Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}
