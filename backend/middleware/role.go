package middleware

import (
	"net/http"

	"smartfarming/dto"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "Unauthorized: role not found in session",
			})
			c.Abort()
			return
		}

		role, ok := roleVal.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "Unauthorized: invalid session parameters",
			})
			c.Abort()
			return
		}

		for _, r := range allowedRoles {
			if r == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, dto.APIResponse{
			Success: false,
			Message: "Access forbidden: you do not have permission to perform this action",
		})
		c.Abort()
	}
}
