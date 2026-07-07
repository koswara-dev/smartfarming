package tests

import (
	"context"
	"testing"

	"smartfarming/dto"
	"smartfarming/service"
)

func TestUserService_CRUD(t *testing.T) {
	repo := newMockUserRepository()
	s := service.NewUserService(repo)

	// Test Create
	createReq := dto.CreateUserRequest{
		Name:     "Operator John",
		Email:    "john@farming.com",
		Password: "johnpassword",
		Role:     "operator",
	}

	userRes, err := s.Create(context.Background(), createReq)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if userRes.Name != "Operator John" || userRes.Role != "operator" {
		t.Errorf("Expected Operator John (operator), got %s (%s)", userRes.Name, userRes.Role)
	}

	// Test GetByID
	fetched, err := s.GetByID(context.Background(), userRes.ID)
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}
	if fetched.ID != userRes.ID {
		t.Errorf("ID mismatch: %v vs %v", fetched.ID, userRes.ID)
	}

	// Test Update
	newName := "Operator John Updated"
	updateReq := dto.UpdateUserRequest{
		Name: &newName,
	}
	updated, err := s.Update(context.Background(), userRes.ID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}
	if updated.Name != "Operator John Updated" {
		t.Errorf("Expected name Operator John Updated, got %s", updated.Name)
	}

	// Test List (Pagination)
	listRes, err := s.List(context.Background(), dto.PaginationRequest{
		Page:  1,
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}
	if listRes.Meta.TotalRecords != 1 {
		t.Errorf("Expected 1 record, got %d", listRes.Meta.TotalRecords)
	}

	// Test Delete
	err = s.Delete(context.Background(), userRes.ID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	_, err = s.GetByID(context.Background(), userRes.ID)
	if err == nil {
		t.Error("Expected error fetching deleted user, got nil")
	}
}
