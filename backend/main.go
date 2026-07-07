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

	// 3. Setup Layers
	userRepo := repository.NewUserRepository(config.DB)

	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	// 4. Initialize Router
	r := routes.SetupRouter(authHandler, userHandler)

	// 5. Start Server
	log.Printf("Server starting on port %s in %s mode...", config.AppConfig.Port, config.AppConfig.Env)
	if err := r.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
