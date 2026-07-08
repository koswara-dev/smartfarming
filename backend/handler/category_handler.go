package handler

import (
	"net/http"

	"smartfarming/dto"
	"smartfarming/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// Create godoc
// @Summary Create a Category
// @Description Creates a new article category (Admin/Operator only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateCategoryRequest true "Create Category request"
// @Success 201 {object} dto.APIResponse{data=dto.CategoryResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 403 {object} dto.APIResponse
// @Router /categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Validation error",
			Errors:  err.Error(),
		})
		return
	}

	res, err := h.categoryService.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "Category created successfully",
		Data:    res,
	})
}

// GetByID godoc
// @Summary Get Category by ID
// @Description Fetches details of a category by its UUID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category UUID"
// @Success 200 {object} dto.APIResponse{data=dto.CategoryResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		// Fallback to fetch by slug if parse fails
		h.GetBySlug(c, idStr)
		return
	}

	res, err := h.categoryService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: "Category not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Fetched category details successfully",
		Data:    res,
	})
}

func (h *CategoryHandler) GetBySlug(c *gin.Context, slug string) {
	res, err := h.categoryService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: "Category not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Fetched category details successfully",
		Data:    res,
	})
}

// Update godoc
// @Summary Update Category
// @Description Updates category details by its UUID (Admin/Operator only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category UUID"
// @Param request body dto.UpdateCategoryRequest true "Update Category request"
// @Success 200 {object} dto.APIResponse{data=dto.CategoryResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 403 {object} dto.APIResponse
// @Router /categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid ID parameter format",
		})
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Validation error",
			Errors:  err.Error(),
		})
		return
	}

	res, err := h.categoryService.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Category updated successfully",
		Data:    res,
	})
}

// Delete godoc
// @Summary Delete Category
// @Description Deletes a category by its UUID (Admin/Operator only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category UUID"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 403 {object} dto.APIResponse
// @Router /categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid ID parameter format",
		})
		return
	}

	err = h.categoryService.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Category deleted successfully",
	})
}

// List godoc
// @Summary List Categories with search and pagination
// @Description Lists categories with pagination and search parameters
// @Tags Categories
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Param search query string false "Search by name or description"
// @Success 200 {object} dto.APIResponse{data=dto.PaginationResponse}
// @Failure 400 {object} dto.APIResponse
// @Router /categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Failed to parse query parameters",
			Errors:  err.Error(),
		})
		return
	}

	res, err := h.categoryService.List(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Categories fetched successfully",
		Data:    res,
	})
}
