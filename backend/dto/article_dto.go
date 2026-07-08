package dto

import "github.com/google/uuid"

type CreateArticleRequest struct {
	Title      string    `json:"title" binding:"required,min=3,max=255"`
	Content    string    `json:"content" binding:"required,min=10"`
	CategoryID uuid.UUID `json:"categoryId" binding:"required"`
}

type UpdateArticleRequest struct {
	Title      *string    `json:"title" binding:"omitempty,min=3,max=255"`
	Content    *string    `json:"content" binding:"omitempty,min=10"`
	CategoryID *uuid.UUID `json:"categoryId" binding:"omitempty"`
}

type ArticleResponse struct {
	ID        uuid.UUID        `json:"id"`
	Title     string           `json:"title"`
	Slug      string           `json:"slug"`
	Content   string           `json:"content"`
	Category  CategoryResponse `json:"category"`
	ImageURL  string           `json:"imageUrl"`
	CreatedBy *uuid.UUID       `json:"createdBy,omitempty"`
	CreatedAt string           `json:"createdAt"`
	UpdatedAt string           `json:"updatedAt"`
}
