package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/rojgarsetu/backend/internal/db"
)

type UploadService struct {
	database *db.PostgresDB
	uploadDir string
}

func NewUploadService(database *db.PostgresDB) *UploadService {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	_ = os.MkdirAll(uploadDir, 0755)
	return &UploadService{database: database, uploadDir: uploadDir}
}

func (s *UploadService) SaveFile(ctx context.Context, userID string, fileType string, fileName string, fileSize int64, reader io.Reader) (*db.FileUpload, error) {
	userUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	// Create unique file name
	ext := filepath.Ext(fileName)
	uniqueName := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	filePath := filepath.Join(s.uploadDir, uniqueName)

	// Save to disk
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, reader); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Store in DB
	fileURL := "/uploads/" + uniqueName // Placeholder URL
	row, err := s.database.Queries.CreateFileUpload(ctx, db.CreateFileUploadParams{
		UserID:       userUID,
		FileType:     fileType,
		FileUrl:      fileURL,
		OriginalName: fileName,
		FileSize:     db.NullInt64(fileSize),
	})
	if err != nil {
		_ = os.Remove(filePath) // Cleanup on DB error
		return nil, err
	}

	return &row, nil
}
