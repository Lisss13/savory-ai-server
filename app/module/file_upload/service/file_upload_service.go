package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"savory-ai-server/app/module/file_upload/payload"
	"savory-ai-server/utils/config"
	"strings"
	"time"
)

type fileUploadService struct {
	cfg *config.Config
}

type FileUploadService interface {
	UploadFile(file *multipart.FileHeader, folder string) (*payload.FileUploadResp, error)
}

func NewFileUploadService(cfg *config.Config) FileUploadService {
	return &fileUploadService{
		cfg: cfg,
	}
}

func (s *fileUploadService) UploadFile(file *multipart.FileHeader, folder string) (*payload.FileUploadResp, error) {
	// Create the upload directory if it doesn't exist
	uploadDir := filepath.Join("storage", "public", folder)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Generate a unique filename to avoid collisions
	filename := generateUniqueFilename(file.Filename)
	filePath := filepath.Join(uploadDir, filename)

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	if _, err = io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// Generate the URL for the file
	// In a production environment, this would be a full URL with domain
	url := fmt.Sprintf("/%s/%s", folder, filename)

	return &payload.FileUploadResp{
		Filename:  filename,
		Size:      file.Size,
		URL:       url,
		CreatedAt: time.Now(),
	}, nil
}

// generateUniqueFilename adds a timestamp to the filename to make it unique
func generateUniqueFilename(filename string) string {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s_%d%s", name, timestamp, ext)
}
