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

func generateTestToken(userID uuid.UUID, role string) string {
	claims := jwt.MapClaims{
		"sub":  userID.String(),
		"role": role,
		"exp":  time.Now().Add(time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("testsecretkeyforunittesting12345"))
	return tokenStr
}

func TestRBAC_SecurityControls(t *testing.T) {
	uRepo := newMockUserRepository()
	catRepo := newMockCategoryRepository()
	artRepo := newMockArticleRepository()

	authSvc := service.NewAuthService(uRepo)
	userSvc := service.NewUserService(uRepo)
	catSvc := service.NewCategoryService(catRepo)
	artSvc := service.NewArticleService(artRepo, catRepo)

	storageSvc := service.NewStorageService(nil, "smartfarming")

	authHand := handler.NewAuthHandler(authSvc)
	userHand := handler.NewUserHandler(userSvc, storageSvc)
	catHand := handler.NewCategoryHandler(catSvc)
	artHand := handler.NewArticleHandler(artSvc, storageSvc)

	router := routes.SetupRouter(authHand, userHand, catHand, artHand)

	userToken := generateTestToken(uuid.New(), "user")
	operatorToken := generateTestToken(uuid.New(), "operator")
	adminToken := generateTestToken(uuid.New(), "admin")

	tests := []struct {
		name           string
		method         string
		url            string
		token          string
		body           string
		expectedStatus int
	}{
		{
			name:           "Standard user cannot create category",
			method:         "POST",
			url:            "/api/v1/categories",
			token:          userToken,
			body:           `{"name":"Tech","description":"Tech desc"}`,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Unauthenticated request to create category fails",
			method:         "POST",
			url:            "/api/v1/categories",
			token:          "",
			body:           `{"name":"Tech","description":"Tech desc"}`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Operator can create category",
			method:         "POST",
			url:            "/api/v1/categories",
			token:          operatorToken,
			body:           `{"name":"Tech","description":"Tech desc"}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Admin can create category",
			method:         "POST",
			url:            "/api/v1/categories",
			token:          adminToken,
			body:           `{"name":"Tech2","description":"Tech desc"}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Standard user cannot create article",
			method:         "POST",
			url:            "/api/v1/articles",
			token:          userToken,
			body:           `{"title":"New article","content":"Content goes here...","categoryId":"` + uuid.New().String() + `"}`,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Public request can read categories list",
			method:         "GET",
			url:            "/api/v1/categories?page=1&limit=10",
			token:          "",
			body:           "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Public request can read articles list",
			method:         "GET",
			url:            "/api/v1/articles?page=1&limit=10",
			token:          "",
			body:           "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unauthenticated request to upload user photo fails",
			method:         "POST",
			url:            "/api/v1/users/photo",
			token:          "",
			body:           "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Standard user cannot upload article image",
			method:         "POST",
			url:            "/api/v1/articles/" + uuid.New().String() + "/image",
			token:          userToken,
			body:           "",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d. Response: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if w.Code == http.StatusForbidden {
				var resp dto.APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				if err != nil {
					t.Errorf("failed to unmarshal JSON error response: %v", err)
				}
				if resp.Success {
					t.Error("expected Success to be false on forbidden error")
				}
			}
		})
	}
}
