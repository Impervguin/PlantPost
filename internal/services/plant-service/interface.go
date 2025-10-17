package plantservice

import (
	"PlantSite/internal/models"
	"PlantSite/internal/models/plant"
	"context"

	"github.com/google/uuid"
)

type PlantServiceContract interface {
	UpdatePlantSpec(ctx context.Context, id uuid.UUID, spec plant.PlantSpecification) error
	DeletePlant(ctx context.Context, id uuid.UUID) error
	UploadPlantPhoto(ctx context.Context, id uuid.UUID, fdata models.FileData, description string) error
	GetPlantCategory(ctx context.Context, name string) (*plant.PlantCategory, error)
	ListCategories(ctx context.Context) ([]plant.PlantCategory, error)
	GetPlant(ctx context.Context, id uuid.UUID) (*GetPlant, error)
	CreatePlant(ctx context.Context, data CreatePlantData, mainPhotoFile models.FileData) error
}
