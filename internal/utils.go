package internal

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func SaveUploadedFile(file *multipart.FileHeader, baseUploadPath string) (string, error) {
	if file == nil {
		return "", fmt.Errorf("file cannot be nil")
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	uniqueFileName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), file.Filename)

	filePath := filepath.Join(baseUploadPath, uniqueFileName)

	if err := os.MkdirAll(baseUploadPath, 0750); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	return filePath, nil
}
