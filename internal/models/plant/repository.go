package plant

import (
	"context"

	"github.com/google/uuid"
)

type PlantRepository interface {
	Create(ctx context.Context, plant *Plant) (*Plant, error)
	Update(ctx context.Context, plantID uuid.UUID, updateFn func(*Plant) (*Plant, error)) (*Plant, error)
	Delete(ctx context.Context, plantID uuid.UUID) error
	Get(ctx context.Context, plantID uuid.UUID) (*Plant, error)
}

type PlantCategoryRepository interface {
	GetCategories(ctx context.Context) ([]PlantCategory, error)
	GetCategory(ctx context.Context, name string) (*PlantCategory, error)
}
