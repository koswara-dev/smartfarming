package middleware

import (
	"net/http"

	"smartfarming/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UserIDORMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authenticatedUserIDVal, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		authenticatedUserID, ok := authenticatedUserIDVal.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "Invalid user session",
			})
			c.Abort()
			return
		}

		roleVal, _ := c.Get("role")
		role, _ := roleVal.(string)

		// Admin has access to all resources
		if role == "admin" {
			c.Next()
			return
		}

		// Extract target user ID from the path parameter
		pathIDStr := c.Param("id")
		if pathIDStr != "" {
			pathID, err := uuid.Parse(pathIDStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, dto.APIResponse{
					Success: false,
					Message: "Invalid ID parameter format",
				})
				c.Abort()
				return
			}

			// Validate if the authenticated user matches the target resource ID
			if authenticatedUserID != pathID {
				c.JSON(http.StatusForbidden, dto.APIResponse{
					Success: false,
					Message: "Access forbidden: you do not have permission to access or modify this resource",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
