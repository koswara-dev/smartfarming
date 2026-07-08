package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	Env           string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	RedisAddr     string
	RedisPassword string
	RedisDB       string
	MongoURI      string
	MongoDBName   string
	JWTSecret     string
	MinioEndpoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioUseSSL     string
	MinioBucketName string
}

var AppConfig *Config

func LoadConfig() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	envFile := ".env." + env
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: error loading %s, falling back to OS environment variables: %v", envFile, err)
	}

	AppConfig = &Config{
		Port:            getEnv("PORT", "8080"),
		Env:             getEnv("ENV", env),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "smartfarming"),
		DBSSLMode:       getEnv("DB_SSLMODE", "disable"),
		RedisAddr:       getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		RedisDB:         getEnv("REDIS_DB", "0"),
		MongoURI:        getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName:     getEnv("MONGO_DB", "smartfarming"),
		JWTSecret:       getEnv("JWT_SECRET", "defaultsecretkeythatisverylong"),
		MinioEndpoint:   getEnv("MINIO_ENDPOINT", "172.30.54.28:9000"),
		MinioAccessKey:  getEnv("MINIO_ACCESS_KEY", "admin"),
		MinioSecretKey:  getEnv("MINIO_SECRET_KEY", "Devjc@1324"),
		MinioUseSSL:     getEnv("MINIO_USE_SSL", "false"),
		MinioBucketName: getEnv("MINIO_BUCKET_NAME", "smartfarming"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
