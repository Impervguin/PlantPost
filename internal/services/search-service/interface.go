package searchservice

import (
	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/post"
	"PlantSite/internal/models/search"
	"context"

	"github.com/google/uuid"
)

type SearchServiceContract interface {
	PostTags(ctx context.Context) ([]string, error)
	PostAuthors(ctx context.Context) ([]*auth.Author, error)
	SearchPosts(ctx context.Context, plSearch *search.PostSearch) ([]*SearchPost, error)
	SearchPlants(ctx context.Context, plSearch *search.PlantSearch) ([]*SearchPlant, error)
	GetPostByID(ctx context.Context, id uuid.UUID) (*post.Post, error)
	GetPost(ctx context.Context, id uuid.UUID) (*GetPost, error)
	GetPlantByID(ctx context.Context, id uuid.UUID) (*GetPlant, error)
}
