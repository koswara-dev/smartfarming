package handler

import (
	"net/http"

	"smartfarming/dto"
	"smartfarming/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Create godoc
// @Summary Create a User
// @Description Creates a new user (Admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateUserRequest true "Create User request"
// @Success 201 {object} dto.APIResponse{data=dto.UserResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 403 {object} dto.APIResponse
// @Router /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	roleVal, _ := c.Get("role")
	if roleVal != "admin" {
		c.JSON(http.StatusForbidden, dto.APIResponse{
			Success: false,
			Message: "Access denied: admin privilege required",
		})
		return
	}

	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Validation error",
			Errors:  err.Error(),
		})
		return
	}

	res, err := h.userService.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "User created successfully",
		Data:    res,
	})
}

// GetByID godoc
// @Summary Get User by ID
// @Description Fetches details of a user by their UUID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User UUID"
// @Success 200 {object} dto.APIResponse{data=dto.UserResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 403 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid ID parameter format",
		})
		return
	}

	res, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Fetched user details successfully",
		Data:    res,
	})
}

// Update godoc
// @Summary Update User
// @Description Updates user details by their UUID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User UUID"
// @Param request body dto.UpdateUserRequest true "Update User request"
// @Success 200 {object} dto.APIResponse{data=dto.UserResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 403 {object} dto.APIResponse
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid ID parameter format",
		})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Validation error",
			Errors:  err.Error(),
		})
		return
	}

	res, err := h.userService.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "User updated successfully",
		Data:    res,
	})
}

// Delete godoc
// @Summary Delete User
// @Description Soft-deletes user record
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User UUID"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 403 {object} dto.APIResponse
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid ID parameter format",
		})
		return
	}

	err = h.userService.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}

// List godoc
// @Summary List Users with search and pagination
// @Description Lists users with pagination and search parameters (Admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Param search query string false "Search by name or email"
// @Success 200 {object} dto.APIResponse{data=dto.PaginationResponse}
// @Failure 401 {object} dto.APIResponse
// @Failure 403 {object} dto.APIResponse
// @Router /users [get]
func (h *UserHandler) List(c *gin.Context) {
	roleVal, _ := c.Get("role")
	if roleVal != "admin" {
		c.JSON(http.StatusForbidden, dto.APIResponse{
			Success: false,
			Message: "Access denied: admin privilege required",
		})
		return
	}

	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Failed to parse query parameters",
			Errors:  err.Error(),
		})
		return
	}

	res, err := h.userService.List(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Users fetched successfully",
		Data:    res,
	})
}
