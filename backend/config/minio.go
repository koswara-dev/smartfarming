package config

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func ConnectMinio() {
	endpoint := AppConfig.MinioEndpoint
	accessKeyID := AppConfig.MinioAccessKey
	secretAccessKey := AppConfig.MinioSecretKey
	useSSL, err := strconv.ParseBool(AppConfig.MinioUseSSL)
	if err != nil {
		useSSL = false
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Printf("Warning: Failed to initialize MinIO client: %v", err)
		return
	}

	MinioClient = client
	log.Println("MinIO client initialized successfully.")

	bucketName := AppConfig.MinioBucketName
	ctx := context.Background()
	exists, err := MinioClient.BucketExists(ctx, bucketName)
	if err != nil {
		log.Printf("Warning: Failed to check if MinIO bucket exists: %v", err)
		return
	}

	if !exists {
		err = MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("Warning: Failed to create MinIO bucket: %v", err)
			return
		}
		log.Printf("MinIO bucket '%s' created successfully.", bucketName)

		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": "*",
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`, bucketName)

		err = MinioClient.SetBucketPolicy(ctx, bucketName, policy)
		if err != nil {
			log.Printf("Warning: Failed to set public policy on MinIO bucket: %v", err)
		} else {
			log.Printf("MinIO bucket '%s' read policy set to public read-only.", bucketName)
		}
	} else {
		log.Printf("MinIO bucket '%s' already exists.", bucketName)
	}
}
