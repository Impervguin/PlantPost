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
	pl, err := s.searchRepo.GetPlantByID(ctx, id)
	if err != nil {
		return nil, Wrap(err)
	}
	mainPhoto, err := s.plantFileRepo.Get(ctx, pl.MainPhotoID())
	if err != nil {
		return nil, Wrap(err)
	}

	photos := make([]GetPlantPhoto, 0)

	err = pl.GetPhotos().Iterate(func(e plant.PlantPhoto) error {
		photoFile, err := s.plantFileRepo.Get(ctx, e.FileID())
		if err != nil {
			return Wrap(err)
		}
		photos = append(photos, GetPlantPhoto{
			ID:          e.ID(),
			File:        *photoFile,
			Description: e.Description(),
		})
		return nil
	})
	if err != nil {
		return nil, Wrap(err)
	}
	return &GetPlant{
		ID:            pl.ID(),
		Name:          pl.GetName(),
		LatinName:     pl.GetLatinName(),
		Description:   pl.GetDescription(),
		MainPhoto:     *mainPhoto,
		Photos:        photos,
		Category:      pl.GetCategory(),
		Specification: pl.GetSpecification(),
		CreatedAt:     pl.CreatedAt(),
	}, nil
}
