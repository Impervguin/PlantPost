package search

import (
	"PlantSite/internal/models/plant"
	"PlantSite/internal/models/post"
	"context"

	"github.com/google/uuid"
)

type SearchRepository interface {
	SearchPosts(ctx context.Context, search *PostSearch) ([]*post.Post, error)
	SearchPlants(ctx context.Context, search *PlantSearch) ([]*plant.Plant, error)
	GetPostByID(ctx context.Context, id uuid.UUID) (*post.Post, error)
	GetPlantByID(ctx context.Context, id uuid.UUID) (*plant.Plant, error)
}
