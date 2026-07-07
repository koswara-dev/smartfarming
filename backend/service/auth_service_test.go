package service

import (
	"context"
	"errors"
	"testing"

	"smartfarming/config"
	"smartfarming/dto"
	"smartfarming/model"

	"github.com/google/uuid"
)

func init() {
	config.AppConfig = &config.Config{
		JWTSecret: "testsecretkeyforunittesting12345",
	}
}

type mockUserRepository struct {
	users  map[uuid.UUID]*model.User
	emails map[string]*model.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:  make(map[uuid.UUID]*model.User),
		emails: make(map[string]*model.User),
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *model.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	m.users[user.ID] = user
	m.emails[user.Email] = user
	return nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	u, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	u, exists := m.emails[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *model.User) error {
	m.users[user.ID] = user
	m.emails[user.Email] = user
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	u, exists := m.users[id]
	if !exists {
		return errors.New("user not found")
	}
	delete(m.users, id)
	delete(m.emails, u.Email)
	return nil
}

func (m *mockUserRepository) List(ctx context.Context, page int, limit int, search string) ([]model.User, int64, error) {
	var list []model.User
	for _, u := range m.users {
		list = append(list, *u)
	}
	return list, int64(len(list)), nil
}

func TestAuthService_Register(t *testing.T) {
	repo := newMockUserRepository()
	s := NewAuthService(repo)

	req := dto.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	res, err := s.Register(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if res.User.Email != "test@example.com" {
		t.Errorf("Expected email to be test@example.com, got %s", res.User.Email)
	}

	if res.Token == "" {
		t.Error("Expected non-empty JWT token")
	}
}

func TestAuthService_Logout(t *testing.T) {
	repo := newMockUserRepository()
	s := NewAuthService(repo)

	req := dto.RegisterRequest{
		Name:     "Logout User",
		Email:    "logout@example.com",
		Password: "password123",
	}

	res, err := s.Register(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = s.Logout(context.Background(), res.Token)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
