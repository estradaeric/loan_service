package utils

import (
	"errors"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
)

var ErrInvalidImageFormat = errors.New("only .jpg or .jpeg images are allowed")

func ReadFileBytes(file multipart.File) ([]byte, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func IsValidImage(fileHeader *multipart.FileHeader) error {
	if fileHeader == nil {
		return nil
	}
	ext := strings.ToLower(fileHeader.Filename)
	if !strings.HasSuffix(ext, ".jpg") && !strings.HasSuffix(ext, ".jpeg") {
		return ErrInvalidImageFormat
	}
	return nil
}

func IsValidPDF(fileHeader *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".pdf" && ext != ".jpg" && ext != ".jpeg" {
		return errors.New("only PDF or JPG/JPEG files are allowed")
	}
	return nil
}

func ValidateJPEG(fileHeader *multipart.FileHeader) error {
	if fileHeader == nil {
		return errors.New("file is required")
	}
	if !strings.HasSuffix(fileHeader.Filename, ".jpg") && !strings.HasSuffix(fileHeader.Filename, ".jpeg") {
		return errors.New("only JPG or JPEG files are allowed")
	}
	return nil
}