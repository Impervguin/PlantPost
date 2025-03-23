package plantservice

import (
	authservice "PlantSite/internal/auth-service"
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

func (s *PlantService) GetPlant(ctx context.Context, id uuid.UUID) (*GetPlant, error) {
	user := authservice.UserFromContext(ctx)
	if user == nil {
		return nil, ErrNotAuthorized
	}
	if !user.HasAuthorRights() {
		return nil, ErrNotAuthor
	}
	plant, err := s.plantrepo.Get(ctx, id)
	if err != nil {
		return nil, Wrap(err)
	}

	mainPhoto, err := s.filerepo.Get(ctx, plant.MainPhotoID())
	if err != nil {
		return nil, Wrap(err)
	}
	photos := make([]GetPlantPhoto, 0, len(plant.GetPhotos()))
	for _, p := range plant.GetPhotos() {
		photoFile, err := s.filerepo.Get(ctx, p.FileID())
		if err != nil {
			return nil, Wrap(err)
		}
		photos = append(photos, GetPlantPhoto{
			ID:          p.ID(),
			File:        *photoFile,
			Description: p.Description(),
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
