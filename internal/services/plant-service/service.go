package plantservice

import (
	"PlantSite/internal/models"
	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/plant"
	authservice "PlantSite/internal/services/auth-service"
	"PlantSite/internal/utils/logs"
	"context"

	"github.com/google/uuid"
)

type PlantService struct {
	plantrepo    plant.PlantRepository
	categoryrepo plant.PlantCategoryRepository
	filerepo     models.FileRepository
	auth         authservice.AuthServiceContract
}

func NewPlantService(repository plant.PlantRepository, crep plant.PlantCategoryRepository, filerepo models.FileRepository, auth authservice.AuthServiceContract) *PlantService {
	if repository == nil {
		panic("nil repository")
	}
	if crep == nil {
		panic("nil category repository")
	}
	if filerepo == nil {
		panic("nil file repository")
	}
	if auth == nil {
		panic("nil auth")
	}
	return &PlantService{plantrepo: repository,
		categoryrepo: crep,
		filerepo:     filerepo,
		auth:         auth,
	}
}

func (s *PlantService) UpdatePlantSpec(ctx context.Context, id uuid.UUID, spec plant.PlantSpecification) error {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return auth.ErrNotAuthorized
	}
	logs.Debugf("PlantService.UpdatePlantSpec: got user=%s", user.ID())
	if !user.HasAuthorRights() {
		return auth.ErrNoAuthorRights
	}
	logs.Debugf("PlantService.UpdatePlantSpec: user %s has author rights", user.ID())
	_, err := s.plantrepo.Update(ctx, id, func(p *plant.Plant) (*plant.Plant, error) {
		err := p.UpdateSpec(spec)
		return p, err
	})
	logs.Debugf("PlantService.UpdatePlantSpec: updated plant %s spec to %s", id, spec)
	return err
}

func (s *PlantService) DeletePlant(ctx context.Context, id uuid.UUID) error {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return auth.ErrNotAuthorized
	}
	logs.Debugf("PlantService.DeletePlant: got user=%s", user.ID())
	if !user.HasAuthorRights() {
		return auth.ErrNoAuthorRights
	}
	logs.Debugf("PlantService.DeletePlant: user %s has author rights", user.ID())
	return s.plantrepo.Delete(ctx, id)
}

func (s *PlantService) UploadPlantPhoto(ctx context.Context, id uuid.UUID, fdata models.FileData, description string) error {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return auth.ErrNotAuthorized
	}
	logs.Debugf("PlantService.UploadPlantPhoto: got user=%s", user.ID())
	if !user.HasAuthorRights() {
		return auth.ErrNoAuthorRights
	}
	logs.Debugf("PlantService.UploadPlantPhoto: user %s has author rights", user.ID())
	file, err := s.filerepo.Upload(ctx, &fdata)
	if err != nil {
		return err
	}
	logs.Debugf("PlantService.UploadPlantPhoto: uploaded photo %s", file.ID)
	_, err = s.plantrepo.Update(ctx, id, func(p *plant.Plant) (*plant.Plant, error) {
		photo, err := plant.NewPlantPhoto(file.ID, description)
		if err != nil {
			return nil, err
		}
		err = p.AddPhoto(photo)
		return p, err
	})
	if err != nil {
		return err
	}
	logs.Debugf("PlantService.UploadPlantPhoto: added photo %s to plant %s", file.ID, id)
	return nil
}

func (s *PlantService) GetPlantCategory(ctx context.Context, name string) (*plant.PlantCategory, error) {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return nil, auth.ErrNotAuthorized
	}
	logs.Debugf("PlantService.GetPlantCategory: got user=%s", user.ID())
	if !user.HasAuthorRights() {
		return nil, auth.ErrNoAuthorRights
	}
	logs.Debugf("PlantService.GetPlantCategory: user %s has author rights", user.ID())
	return s.categoryrepo.GetCategory(ctx, name)
}

func (s *PlantService) ListCategories(ctx context.Context) ([]plant.PlantCategory, error) {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return nil, auth.ErrNotAuthorized
	}
	logs.Debugf("PlantService.ListCategories: got user=%s", user.ID())
	if !user.HasAuthorRights() {
		return nil, auth.ErrNoAuthorRights
	}
	logs.Debugf("PlantService.ListCategories: user %s has author rights", user.ID())
	return s.categoryrepo.GetCategories(ctx)
}
