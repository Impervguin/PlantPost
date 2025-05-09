package plantservice

import (
	"PlantSite/internal/models"
	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/plant"
	"context"
	"fmt"
)

type CreatePlantData struct {
	Name        string
	LatinName   string
	Description string
	Category    string
	Spec        plant.PlantSpecification
}

func (s *PlantService) CreatePlant(ctx context.Context, data CreatePlantData, mainPhotoFile models.FileData) error {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return auth.ErrNotAuthorized
	}
	if !user.HasAuthorRights() {
		return auth.ErrNoAuthorRights
	}

	_, err := s.categoryrepo.GetCategory(ctx, data.Category)
	if err != nil {
		return Wrap(err)
	}

	f, err := s.filerepo.Upload(ctx, &mainPhotoFile)
	if err != nil {
		return fmt.Errorf("failed to upload main photo: %w", err)
	}
	plant, err := plant.NewPlant(data.Name, data.LatinName, data.Description, f.ID, *plant.NewPlantPhotos(), data.Category, data.Spec)
	if err != nil {
		return Wrap(err)
	}
	_, err = s.plantrepo.Create(ctx, plant)
	if err != nil {
		return Wrap(err)
	}

	return nil
}
