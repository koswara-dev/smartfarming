package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"smartfarming/config"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type StorageService interface {
	UploadFile(ctx context.Context, fileReader io.Reader, fileSize int64, contentType string, folder string) (string, error)
}

type storageService struct {
	minioClient *minio.Client
	bucketName  string
}

func NewStorageService(minioClient *minio.Client, bucketName string) StorageService {
	return &storageService{
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}

func (s *storageService) UploadFile(ctx context.Context, fileReader io.Reader, fileSize int64, contentType string, folder string) (string, error) {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
		"image/gif":  true,
	}
	if !allowedTypes[contentType] {
		return "", errors.New("unsupported file format. Only JPEG, PNG, WEBP, and GIF are allowed")
	}

	if fileSize > 5*1024*1024 {
		return "", errors.New("file size exceeds maximum limit of 5MB")
	}

	ext := ".jpg"
	if strings.Contains(contentType, "/") {
		ext = "." + strings.Split(contentType, "/")[1]
	}
	if ext == ".jpeg" {
		ext = ".jpg"
	}

	uniqueName := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	if s.minioClient == nil {
		log.Printf("[Mock Storage] Uploaded %s size=%d type=%s to folder=%s", uniqueName, fileSize, contentType, folder)
		return fmt.Sprintf("http://172.30.54.28:9000/smartfarming/%s", uniqueName), nil
	}

	_, err := s.minioClient.PutObject(ctx, s.bucketName, uniqueName, fileReader, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object: %v", err)
	}

	url := fmt.Sprintf("http://%s/%s/%s", config.AppConfig.MinioEndpoint, s.bucketName, uniqueName)
	return url, nil
}
