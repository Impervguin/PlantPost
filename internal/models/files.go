package models

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
)

var (
	ErrFileNotFound = errors.New("file not found")
)

// File inside system
type File struct {
	ID        uuid.UUID
	Name      string
	URL       string // To get as static file
	CreatedAt time.Time
}

func CreateFile(id uuid.UUID, name string, url string, createdAt time.Time) (*File, error) {
	if id == uuid.Nil {
		id = uuid.New()
	}
	return &File{
		ID:        id,
		Name:      name,
		URL:       url,
		CreatedAt: createdAt,
	}, nil
}

func NewFile(name string) (*File, error) {
	if name == "" {
		return nil, fmt.Errorf("name must not be empty")
	}
	id := uuid.New()
	createdAt := time.Now()
	return CreateFile(id, name, id.String(), createdAt)
}

// FileData is for uploading/downloading files to the system
type FileData struct {
	Name        string
	Reader      io.Reader
	ContentType string
}

func NewFileData(name string, reader io.Reader, contentType string) (*FileData, error) {
	return &FileData{
		Name:        name,
		Reader:      reader,
		ContentType: contentType,
	}, nil
}

type FileRepository interface {
	Upload(ctx context.Context, data *FileData) (*File, error)
	Download(ctx context.Context, fileID uuid.UUID) (*FileData, error)
	Delete(ctx context.Context, fileID uuid.UUID) error
	Update(ctx context.Context, fileID uuid.UUID, data *FileData) (*File, error)
	Get(ctx context.Context, fileID uuid.UUID) (*File, error)
}
