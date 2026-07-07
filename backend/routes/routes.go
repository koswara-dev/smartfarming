package routes

import (
	"smartfarming/config"
	"smartfarming/handler"
	"smartfarming/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "smartfarming/docs" // Loaded side-effect imports for auto-generated Swagger docs
)

func SetupRouter(authHandler *handler.AuthHandler, userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()

	// Allow all CORS in development mode
	if config.AppConfig.Env == "development" {
		r.Use(middleware.CORSMiddleware())
	}

	// Swagger API docs route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		// Public Auth Endpoints
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected Auth Endpoints
		authProtected := api.Group("/auth")
		authProtected.Use(middleware.AuthMiddleware())
		{
			authProtected.GET("/me", authHandler.Me)
			authProtected.POST("/logout", authHandler.Logout)
		}

		// User Management Endpoints (All require authentication)
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			// Creating/Listing all users is restricted to Admins in handler logic
			users.POST("", userHandler.Create)
			users.GET("", userHandler.List)

			// IDOR Protected routes: users can view/edit/delete their own record (or Admin)
			usersID := users.Group("/:id")
			usersID.Use(middleware.UserIDORMiddleware())
			{
				usersID.GET("", userHandler.GetByID)
				usersID.PUT("", userHandler.Update)
				usersID.DELETE("", userHandler.Delete)
			}
		}
	}

	return r
}
