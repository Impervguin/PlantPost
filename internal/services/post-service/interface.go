package postservice

import (
	"PlantSite/internal/models"
	"PlantSite/internal/models/post"
	"context"

	"github.com/google/uuid"
)

type PostServiceContract interface {
	CreatePost(ctx context.Context, data CreatePostTextData, files []models.FileData) (*post.Post, error)
	GetPost(ctx context.Context, id uuid.UUID) (*GetPost, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdatePost(ctx context.Context, id uuid.UUID, data UpdatePostTextData) (*post.Post, error)
}
