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

type mockCategoryRepository struct {
	categories map[uuid.UUID]*model.Category
	slugs      map[string]*model.Category
}

func newMockCategoryRepository() *mockCategoryRepository {
	return &mockCategoryRepository{
		categories: make(map[uuid.UUID]*model.Category),
		slugs:      make(map[string]*model.Category),
	}
}

func (m *mockCategoryRepository) Create(ctx context.Context, category *model.Category) error {
	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}
	m.categories[category.ID] = category
	m.slugs[category.Slug] = category
	return nil
}

func (m *mockCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	c, exists := m.categories[id]
	if !exists {
		return nil, errors.New("category not found")
	}
	return c, nil
}

func (m *mockCategoryRepository) GetBySlug(ctx context.Context, slug string) (*model.Category, error) {
	c, exists := m.slugs[slug]
	if !exists {
		return nil, errors.New("category not found")
	}
	return c, nil
}

func (m *mockCategoryRepository) Update(ctx context.Context, category *model.Category) error {
	m.categories[category.ID] = category
	m.slugs[category.Slug] = category
	return nil
}

func (m *mockCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	c, exists := m.categories[id]
	if !exists {
		return errors.New("category not found")
	}
	delete(m.categories, id)
	delete(m.slugs, c.Slug)
	return nil
}

func (m *mockCategoryRepository) List(ctx context.Context, page int, limit int, search string) ([]model.Category, int64, error) {
	var list []model.Category
	for _, c := range m.categories {
		list = append(list, *c)
	}
	return list, int64(len(list)), nil
}

func TestCategoryService_CRUD(t *testing.T) {
	repo := newMockCategoryRepository()
	s := service.NewCategoryService(repo)

	createReq := dto.CreateCategoryRequest{
		Name:        "Hydroponics Guide",
		Description: "Articles about Hydroponics systems",
	}

	res, err := s.Create(context.Background(), createReq)
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	if res.Name != "Hydroponics Guide" || res.Slug != "hydroponics-guide" {
		t.Errorf("Expected name and slug match, got %s (%s)", res.Name, res.Slug)
	}

	// Update
	newName := "Advanced Hydroponics"
	updateReq := dto.UpdateCategoryRequest{
		Name: &newName,
	}

	resUpdated, err := s.Update(context.Background(), res.ID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update category: %v", err)
	}

	if resUpdated.Name != "Advanced Hydroponics" || resUpdated.Slug != "advanced-hydroponics" {
		t.Errorf("Expected updated name/slug, got %s (%s)", resUpdated.Name, resUpdated.Slug)
	}
}
