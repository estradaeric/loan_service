package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// SaveUploadedFile saves the uploaded file to the specified subdirectory and returns the relative path
func SaveUploadedFile(file multipart.File, header *multipart.FileHeader, subdir string) (string, error) {
	defer file.Close()

	// Ensure subdir exists
	fullDir := filepath.Join("/mnt/data/uploads", subdir)
	if err := os.MkdirAll(fullDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	fullPath := filepath.Join(fullDir, filename)

	// Create destination file
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy contents
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return relative path for storage
	return filepath.Join("uploads", subdir, filename), nil
}