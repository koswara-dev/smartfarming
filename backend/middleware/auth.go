package middleware

import (
	"context"
	"net/http"
	"strings"

	"smartfarming/config"
	"smartfarming/dto"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "Authorization header is required",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "Authorization header format must be Bearer {token}",
			})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		// Check if token is blacklisted in Redis
		if config.RedisClient != nil {
			blacklistKey := "blacklist:" + tokenStr
			val, _ := config.RedisClient.Get(c.Request.Context(), blacklistKey).Result()
			if val == "1" {
				c.JSON(http.StatusUnauthorized, dto.APIResponse{
					Success: false,
					Message: "Token has been invalidated (logged out)",
				})
				c.Abort()
				return
			}
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		sub, err := claims.GetSubject()
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "Invalid token claims",
			})
			c.Abort()
			return
		}

		userID, err := uuid.Parse(sub)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "Invalid user ID in token",
			})
			c.Abort()
			return
		}

		var role string
		if roleVal, ok := claims["role"].(string); ok {
			role = roleVal
		}

		c.Set("userID", userID)
		c.Set("role", role)

		// Inject userID into Go context for GORM hooks
		ctx := context.WithValue(c.Request.Context(), "userID", userID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
