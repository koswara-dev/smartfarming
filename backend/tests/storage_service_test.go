package tests

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"smartfarming/service"
)

func TestStorageService_Validation(t *testing.T) {
	s := service.NewStorageService(nil, "smartfarming")

	t.Run("Valid JPEG Upload", func(t *testing.T) {
		content := []byte("fakejpegcontent")
		reader := bytes.NewReader(content)
		url, err := s.UploadFile(context.Background(), reader, int64(len(content)), "image/jpeg", "photos")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !strings.HasPrefix(url, "http://172.30.54.28:9000/smartfarming/photos/") || !strings.HasSuffix(url, ".jpg") {
			t.Errorf("Expected URL structure with folder and suffix .jpg, got %s", url)
		}
	})

	t.Run("Oversized File Fails", func(t *testing.T) {
		content := make([]byte, 6*1024*1024)
		reader := bytes.NewReader(content)
		_, err := s.UploadFile(context.Background(), reader, int64(len(content)), "image/png", "photos")
		if err == nil {
			t.Fatal("Expected error for oversized file, got nil")
		}
		if !strings.Contains(err.Error(), "exceeds maximum limit") {
			t.Errorf("Expected size limit error message, got: %v", err)
		}
	})

	t.Run("Unsupported Format Fails", func(t *testing.T) {
		content := []byte("fakepdfcontent")
		reader := bytes.NewReader(content)
		_, err := s.UploadFile(context.Background(), reader, int64(len(content)), "application/pdf", "photos")
		if err == nil {
			t.Fatal("Expected error for unsupported MIME format, got nil")
		}
		if !strings.Contains(err.Error(), "unsupported file format") {
			t.Errorf("Expected unsupported format error, got: %v", err)
		}
	})
}
