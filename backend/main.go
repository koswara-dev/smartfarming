package main

import (
	"log"
	"smartfarming/config"
	"smartfarming/handler"
	"smartfarming/repository"
	"smartfarming/routes"
	"smartfarming/service"
)

// @title SmartFarming API
// @version 1.0
// @description Backend API for SmartFarming fullstack application using Gin, GORM, PostgreSQL, MongoDB, and Redis.
// @host localhost:8081
// @BasePath /api/v1

// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer " followed by your JWT token to authenticate.
func main() {
	// 1. Load Configurations
	config.LoadConfig()

	// Initialize File Logger
	config.InitLogger()

	// 2. Connect to Databases
	config.ConnectPostgres()
	config.ConnectMongoDB()
	config.ConnectRedis()
	config.ConnectMinio()

	// 3. Setup Layers
	userRepo := repository.NewUserRepository(config.DB)
	categoryRepo := repository.NewCategoryRepository(config.DB)
	articleRepo := repository.NewArticleRepository(config.DB)

	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	articleService := service.NewArticleService(articleRepo, categoryRepo)
	storageService := service.NewStorageService(config.MinioClient, config.AppConfig.MinioBucketName)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService, storageService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	articleHandler := handler.NewArticleHandler(articleService, storageService)

	// 4. Initialize Router
	r := routes.SetupRouter(authHandler, userHandler, categoryHandler, articleHandler)

	// 5. Start Server
	log.Printf("Server starting on port %s in %s mode...", config.AppConfig.Port, config.AppConfig.Env)
	if err := r.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
