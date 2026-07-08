package tests

import (
	"context"
	"errors"
	"testing"

	"smartfarming/dto"
	"smartfarming/model"
	"smartfarming/service"

	"github.com/google/uuid"
)

type mockArticleRepository struct {
	articles map[uuid.UUID]*model.Article
	slugs    map[string]*model.Article
}

func newMockArticleRepository() *mockArticleRepository {
	return &mockArticleRepository{
		articles: make(map[uuid.UUID]*model.Article),
		slugs:    make(map[string]*model.Article),
	}
}

func (m *mockArticleRepository) Create(ctx context.Context, article *model.Article) error {
	if article.ID == uuid.Nil {
		article.ID = uuid.New()
	}
	m.articles[article.ID] = article
	m.slugs[article.Slug] = article
	return nil
}

func (m *mockArticleRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Article, error) {
	a, exists := m.articles[id]
	if !exists {
		return nil, errors.New("article not found")
	}
	return a, nil
}

func (m *mockArticleRepository) GetBySlug(ctx context.Context, slug string) (*model.Article, error) {
	a, exists := m.slugs[slug]
	if !exists {
		return nil, errors.New("article not found")
	}
	return a, nil
}

func (m *mockArticleRepository) Update(ctx context.Context, article *model.Article) error {
	m.articles[article.ID] = article
	m.slugs[article.Slug] = article
	return nil
}

func (m *mockArticleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	a, exists := m.articles[id]
	if !exists {
		return errors.New("article not found")
	}
	delete(m.articles, id)
	delete(m.slugs, a.Slug)
	return nil
}

func (m *mockArticleRepository) List(ctx context.Context, page int, limit int, search string, categoryID *uuid.UUID) ([]model.Article, int64, error) {
	var list []model.Article
	for _, a := range m.articles {
		list = append(list, *a)
	}
	return list, int64(len(list)), nil
}

func TestArticleService_CRUD(t *testing.T) {
	catRepo := newMockCategoryRepository()
	artRepo := newMockArticleRepository()

	// Seed Category
	cat := &model.Category{
		Name: "Farming Tech",
		Slug: "farming-tech",
	}
	catRepo.Create(context.Background(), cat)

	s := service.NewArticleService(artRepo, catRepo)

	userID := uuid.New()
	createReq := dto.CreateArticleRequest{
		Title:      "How to grow mint",
		Content:    "Mint needs a lot of water and partial shade...",
		CategoryID: cat.ID,
	}

	res, err := s.Create(context.Background(), createReq, userID)
	if err != nil {
		t.Fatalf("Failed to create article: %v", err)
	}

	if res.Title != "How to grow mint" || res.Slug != "how-to-grow-mint" {
		t.Errorf("Expected title and slug match, got %s (%s)", res.Title, res.Slug)
	}

	if res.CreatedBy == nil || *res.CreatedBy != userID {
		t.Errorf("Expected CreatedBy to match %s, got %v", userID, res.CreatedBy)
	}

	// Update
	newContent := "Mint needs a lot of water and partial shade. It spreads quickly."
	updateReq := dto.UpdateArticleRequest{
		Content: &newContent,
	}

	resUpdated, err := s.Update(context.Background(), res.ID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update article: %v", err)
	}

	if resUpdated.Content != newContent {
		t.Errorf("Expected content updated, got %s", resUpdated.Content)
	}
}
