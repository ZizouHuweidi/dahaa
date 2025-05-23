package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// ImageStorage handles image file operations
type ImageStorage struct {
	basePath string
}

// NewImageStorage creates a new image storage instance
func NewImageStorage(basePath string) (*ImageStorage, error) {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}

	return &ImageStorage{
		basePath: basePath,
	}, nil
}

// SaveImage saves an uploaded image file
func (s *ImageStorage) SaveImage(file *multipart.FileHeader) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	path := filepath.Join(s.basePath, filename)

	// Open source file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file contents
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file contents: %w", err)
	}

	return filename, nil
}

// GetImagePath returns the full path to an image
func (s *ImageStorage) GetImagePath(filename string) string {
	return filepath.Join(s.basePath, filename)
}

// DeleteImage deletes an image file
func (s *ImageStorage) DeleteImage(filename string) error {
	path := filepath.Join(s.basePath, filename)
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}
	return nil
}

// ValidateImage validates an uploaded image
func (s *ImageStorage) ValidateImage(file *multipart.FileHeader) error {
	// Check file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return fmt.Errorf("file too large: maximum size is 5MB")
	}

	// Check file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}

	if !allowedExts[ext] {
		return fmt.Errorf("invalid file type: only jpg, jpeg, png, and gif are allowed")
	}

	return nil
}
