package plantservice

import (
	authservice "PlantSite/internal/auth-service"
	"PlantSite/internal/models"
	"PlantSite/internal/models/plant"
	"context"
)

type CreatePlantData struct {
	Name        string
	LatinName   string
	Description string
	Category    string
	Spec        plant.PlantSpecification
}

func (s *PlantService) CreatePlant(ctx context.Context, data CreatePlantData, mainPhotoFile models.FileData) error {
	user := authservice.UserFromContext(ctx)
	if user == nil {
		return ErrNotAuthorized
	}
	if !user.HasAuthorRights() {
		return ErrNotAuthor
	}

	_, err := s.categoryrepo.GetCategory(ctx, data.Category)
	if err != nil {
		return Wrap(err)
	}

	f, err := s.filerepo.Upload(ctx, &mainPhotoFile)
	if err != nil {
		return Wrap(err)
	}
	plant, err := plant.NewPlant(data.Name, data.LatinName, data.Description, f.ID, []plant.PlantPhoto{}, data.Category, data.Spec)
	if err != nil {
		return Wrap(err)
	}
	_, err = s.plantrepo.Create(ctx, plant)
	if err != nil {
		return Wrap(err)
	}

	return nil
}
