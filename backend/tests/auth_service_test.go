package tests

import (
	"context"
	"errors"
	"testing"

	"smartfarming/config"
	"smartfarming/dto"
	"smartfarming/model"
	"smartfarming/service"

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
	s := service.NewAuthService(repo)

	req := dto.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	res, err := s.Register(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if res.OTP == "" {
		t.Fatal("Expected generated OTP code, got empty")
	}

	// Verify incorrect OTP
	verifyWrongReq := dto.VerifyOTPRequest{
		Email: "test@example.com",
		Code:  "999999",
	}
	_, err = s.VerifyOTP(context.Background(), verifyWrongReq)
	if err == nil {
		t.Error("Expected error verifying with incorrect OTP code, got nil")
	}

	// Verify correct OTP
	verifyReq := dto.VerifyOTPRequest{
		Email: "test@example.com",
		Code:  res.OTP,
	}
	authRes, err := s.VerifyOTP(context.Background(), verifyReq)
	if err != nil {
		t.Fatalf("Expected no verification error, got %v", err)
	}

	if authRes.User.Email != "test@example.com" {
		t.Errorf("Expected email to be test@example.com, got %s", authRes.User.Email)
	}

	if authRes.Token == "" {
		t.Error("Expected non-empty JWT token")
	}
}

func TestAuthService_Logout(t *testing.T) {
	repo := newMockUserRepository()
	s := service.NewAuthService(repo)

	req := dto.RegisterRequest{
		Name:     "Logout User",
		Email:    "logout@example.com",
		Password: "password123",
	}

	res, err := s.Register(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	verifyReq := dto.VerifyOTPRequest{
		Email: "logout@example.com",
		Code:  res.OTP,
	}
	authRes, err := s.VerifyOTP(context.Background(), verifyReq)
	if err != nil {
		t.Fatalf("Expected no verification error, got %v", err)
	}

	err = s.Logout(context.Background(), authRes.Token)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
