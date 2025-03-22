package models

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
)

// File inside system
type File struct {
	ID        uuid.UUID
	Name      string
	URL       string // To get as static file
	CreatedAt time.Time
}

// FileData is for uploading/downloading files to the system
type FileData struct {
	Name   string
	Reader io.Reader
}

type FileRepository interface {
	Upload(ctx context.Context, data *FileData) (*File, error)
	Download(ctx context.Context, fileID uuid.UUID) (*FileData, error)
	Delete(ctx context.Context, fileID uuid.UUID) error
	Update(ctx context.Context, fileID uuid.UUID, data *FileData) (*File, error)
	Get(ctx context.Context, fileID uuid.UUID) (*File, error)
}
