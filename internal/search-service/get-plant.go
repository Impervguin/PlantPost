package searchservice

import (
	"PlantSite/internal/models"
	"PlantSite/internal/models/plant"
	"context"
	"time"

	"github.com/google/uuid"
)

type GetPlant struct {
	ID            uuid.UUID
	Name          string
	LatinName     string
	Description   string
	MainPhoto     models.File
	Photos        []GetPlantPhoto
	Category      string
	Specification plant.PlantSpecification
	CreatedAt     time.Time
}

type GetPlantPhoto struct {
	ID          uuid.UUID
	File        models.File
	Description string
}

func (s *SearchService) GetPlantByID(ctx context.Context, id uuid.UUID) (*GetPlant, error) {
	plant, err := s.searchRepo.GetPlantByID(ctx, id)
	if err != nil {
		return nil, Wrap(err)
	}
	mainPhoto, err := s.plantFileRepo.Get(ctx, plant.MainPhotoID())
	if err != nil {
		return nil, Wrap(err)
	}

	photos := make([]GetPlantPhoto, 0)
	for _, photo := range plant.GetPhotos() {
		file, err := s.plantFileRepo.Get(ctx, photo.FileID())
		if err != nil {
			return nil, Wrap(err)
		}
		photos = append(photos, GetPlantPhoto{
			ID:          photo.ID(),
			File:        *file,
			Description: photo.Description(),
		})
	}

	return &GetPlant{
		ID:            plant.ID(),
		Name:          plant.GetName(),
		LatinName:     plant.GetLatinName(),
		Description:   plant.GetDescription(),
		MainPhoto:     *mainPhoto,
		Photos:        photos,
		Category:      plant.GetCategory(),
		Specification: plant.GetSpecification(),
		CreatedAt:     plant.CreatedAt(),
	}, nil
}
