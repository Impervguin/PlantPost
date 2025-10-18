package plantservice

import (
	"PlantSite/internal/models"
	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/plant"
	"PlantSite/internal/utils/logs"
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
	logs.Debugf("PlantService.CreatePlant: got user=%s", user.ID())
	if !user.HasAuthorRights() {
		return auth.ErrNoAuthorRights
	}
	logs.Debugf("PlantService.CreatePlant: user %s has author rights", user.ID())

	_, err := s.categoryrepo.GetCategory(ctx, data.Category)
	if err != nil {
		return Wrap(err)
	}
	logs.Debugf("PlantService.CreatePlant: got category %s", data.Category)

	f, err := s.filerepo.Upload(ctx, &mainPhotoFile)
	if err != nil {
		return fmt.Errorf("failed to upload main photo: %w", err)
	}
	logs.Debugf("PlantService.CreatePlant: uploaded main photo %s", f.ID)
	plant, err := plant.NewPlant(data.Name, data.LatinName, data.Description, f.ID, *plant.NewPlantPhotos(), data.Category, data.Spec)
	if err != nil {
		return Wrap(err)
	}
	logs.Debugf("PlantService.CreatePlant: created plant %s", plant.ID())
	_, err = s.plantrepo.Create(ctx, plant)
	if err != nil {
		return Wrap(err)
	}
	logs.Debugf("PlantService.CreatePlant: saved plant %s", plant.ID())

	return nil
}
