package handler

import (
	"net/http"

	"smartfarming/dto"
	"smartfarming/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ArticleHandler struct {
	articleService service.ArticleService
	storageService service.StorageService
}

func NewArticleHandler(articleService service.ArticleService, storageService service.StorageService) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
		storageService: storageService,
	}
}

// Create godoc
// @Summary Create an Article
// @Description Creates a new article (Admin/Operator only)
// @Tags Articles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateArticleRequest true "Create Article request"
// @Success 201 {object} dto.APIResponse{data=dto.ArticleResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 403 {object} dto.APIResponse
// @Router /articles [post]
func (h *ArticleHandler) Create(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "Unauthorized: user ID not found in session",
		})
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "Unauthorized: invalid user session",
		})
		return
	}

	var req dto.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Validation error",
			Errors:  err.Error(),
		})
		return
	}

	res, err := h.articleService.Create(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "Article created successfully",
		Data:    res,
	})
}

// GetByID godoc
// @Summary Get Article by ID or Slug
// @Description Fetches details of an article by its UUID or Slug
// @Tags Articles
// @Accept json
// @Produce json
// @Param id path string true "Article UUID or Slug"
// @Success 200 {object} dto.APIResponse{data=dto.ArticleResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /articles/{id} [get]
func (h *ArticleHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		// Fallback to fetch by slug if parse fails
		h.GetBySlug(c, idStr)
		return
	}

	res, err := h.articleService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: "Article not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Fetched article details successfully",
		Data:    res,
	})
}

func (h *ArticleHandler) GetBySlug(c *gin.Context, slug string) {
	res, err := h.articleService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: "Article not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Fetched article details successfully",
		Data:    res,
	})
}

// Update godoc
// @Summary Update Article
// @Description Updates article details by its UUID (Admin/Operator only)
// @Tags Articles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Article UUID"
// @Param request body dto.UpdateArticleRequest true "Update Article request"
// @Success 200 {object} dto.APIResponse{data=dto.ArticleResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 403 {object} dto.APIResponse
// @Router /articles/{id} [put]
func (h *ArticleHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid ID parameter format",
		})
		return
	}

	var req dto.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Validation error",
			Errors:  err.Error(),
		})
		return
	}

	res, err := h.articleService.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Article updated successfully",
		Data:    res,
	})
}

// Delete godoc
// @Summary Delete Article
// @Description Deletes an article by its UUID (Admin/Operator only)
// @Tags Articles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Article UUID"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 403 {object} dto.APIResponse
// @Router /articles/{id} [delete]
func (h *ArticleHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid ID parameter format",
		})
		return
	}

	err = h.articleService.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Article deleted successfully",
	})
}

// List godoc
// @Summary List Articles with search, pagination and filtering
// @Description Lists articles with pagination, search by title/content, and filter by categoryId
// @Tags Articles
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Param search query string false "Search by title or content"
// @Param categoryId query string false "Filter by Category UUID"
// @Success 200 {object} dto.APIResponse{data=dto.PaginationResponse}
// @Failure 400 {object} dto.APIResponse
// @Router /articles [get]
func (h *ArticleHandler) List(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Failed to parse query parameters",
			Errors:  err.Error(),
		})
		return
	}

	var categoryID *uuid.UUID
	categoryIDStr := c.Query("categoryId")
	if categoryIDStr != "" {
		parsed, err := uuid.Parse(categoryIDStr)
		if err == nil {
			categoryID = &parsed
		}
	}

	res, err := h.articleService.List(c.Request.Context(), req, categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Articles fetched successfully",
		Data:    res,
	})
}

// UploadImage godoc
// @Summary Upload Article Image
// @Description Uploads an image asset for a specific article (Admin/Operator only)
// @Tags Articles
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "Article UUID"
// @Param image formData file true "Image file"
// @Success 200 {object} dto.APIResponse{data=dto.ArticleResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 403 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /articles/{id}/image [post]
func (h *ArticleHandler) UploadImage(c *gin.Context) {
	idStr := c.Param("id")
	articleID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid ID parameter format",
		})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Failed to parse file upload: 'image' field is missing or invalid",
		})
		return
	}
	defer file.Close()

	url, err := h.storageService.UploadFile(c.Request.Context(), file, header.Size, header.Header.Get("Content-Type"), "articles")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	res, err := h.articleService.UpdateImage(c.Request.Context(), articleID, url)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Article image uploaded successfully",
		Data:    res,
	})
}
