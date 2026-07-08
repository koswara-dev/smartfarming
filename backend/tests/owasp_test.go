package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"smartfarming/dto"
	"smartfarming/handler"
	"smartfarming/routes"
	"smartfarming/service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestOWASP_SecurityControls(t *testing.T) {
	uRepo := newMockUserRepository()
	catRepo := newMockCategoryRepository()
	artRepo := newMockArticleRepository()

	authSvc := service.NewAuthService(uRepo)
	userSvc := service.NewUserService(uRepo)
	catSvc := service.NewCategoryService(catRepo)
	artSvc := service.NewArticleService(artRepo, catRepo)

	authHand := handler.NewAuthHandler(authSvc)
	userHand := handler.NewUserHandler(userSvc, nil)
	catHand := handler.NewCategoryHandler(catSvc)
	artHand := handler.NewArticleHandler(artSvc, nil)

	router := routes.SetupRouter(authHand, userHand, catHand, artHand)

	t.Run("A01:2021 - IDOR protection on User update", func(t *testing.T) {
		userID := uuid.New()
		anotherUserID := uuid.New()
		token := generateTestToken(userID, "user")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/users/"+anotherUserID.String(), strings.NewReader(`{"name":"New name"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("expected 403 Forbidden for IDOR access, got %d", w.Code)
		}
	})

	t.Run("A03:2021 - SQL Injection prevention in User List search parameter", func(t *testing.T) {
		token := generateTestToken(uuid.New(), "admin")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users?search='%20OR%20'1'='1", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200 OK for literal SQLi string search, got %d", w.Code)
		}
	})

	t.Run("A07:2021 - Reject expired JWT tokens", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub":  uuid.New().String(),
			"role": "user",
			"exp":  time.Now().Add(-time.Hour).Unix(),
			"iat":  time.Now().Add(-2 * time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		expiredTokenStr, _ := token.SignedString([]byte("testsecretkeyforunittesting12345"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/auth/me", nil)
		req.Header.Set("Authorization", "Bearer "+expiredTokenStr)

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401 Unauthorized for expired JWT, got %d", w.Code)
		}
	})

	t.Run("A08:2021 - Input validation bindings check", func(t *testing.T) {
		token := generateTestToken(uuid.New(), "admin")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/categories", strings.NewReader(`{"name":"A","description":"Short name"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request for short category name, got %d", w.Code)
		}
	})

	t.Run("A03:2021 - XSS Payload sanitization check on Category creation", func(t *testing.T) {
		token := generateTestToken(uuid.New(), "admin")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/categories", strings.NewReader(`{"name":"<script>alert(1)</script> Category","description":"<img src=x onerror=alert(2)>"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created for valid category request, got %d", w.Code)
		}

		var resp dto.APIResponse
		json.Unmarshal(w.Body.Bytes(), &resp)

		dataMap, ok := resp.Data.(map[string]interface{})
		if !ok {
			t.Fatal("failed to typecast response data map")
		}

		escapedName := dataMap["name"].(string)
		escapedDesc := dataMap["description"].(string)

		if strings.Contains(escapedName, "<script>") || !strings.Contains(escapedName, "&lt;script&gt;") {
			t.Errorf("XSS payload in Name was not properly escaped: %s", escapedName)
		}

		if strings.Contains(escapedDesc, "<img") || !strings.Contains(escapedDesc, "&lt;img") {
			t.Errorf("XSS payload in Description was not properly escaped: %s", escapedDesc)
		}
	})
}
