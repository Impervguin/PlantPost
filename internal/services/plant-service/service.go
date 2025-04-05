package plantservice

import (
	"PlantSite/internal/models"
	"PlantSite/internal/models/plant"
	authservice "PlantSite/internal/services/auth-service"
	"context"

	"github.com/google/uuid"
)

type PlantService struct {
	plantrepo    plant.PlantRepository
	categoryrepo plant.PlantCategoryRepository
	filerepo     models.FileRepository
}

func NewPlantService(repository plant.PlantRepository, crep plant.PlantCategoryRepository, filerepo models.FileRepository) *PlantService {
	if repository == nil {
		panic("nil repository")
	}
	if crep == nil {
		panic("nil category repository")
	}
	if filerepo == nil {
		panic("nil file repository")
	}
	return &PlantService{plantrepo: repository,
		categoryrepo: crep,
		filerepo:     filerepo,
	}
}

func (s *PlantService) UpdatePlantSpec(ctx context.Context, id uuid.UUID, spec plant.PlantSpecification) error {

	user := authservice.UserFromContext(ctx)
	if user == nil {
		return ErrNotAuthorized
	}
	if !user.HasAuthorRights() {
		return ErrNotAuthor
	}
	_, err := s.plantrepo.Update(ctx, id, func(p *plant.Plant) (*plant.Plant, error) {
		err := p.UpdateSpec(spec)
		return p, err
	})
	return err
}

func (s *PlantService) DeletePlant(ctx context.Context, id uuid.UUID) error {
	user := authservice.UserFromContext(ctx)
	if user == nil {
		return ErrNotAuthorized
	}
	if !user.HasAuthorRights() {
		return ErrNotAuthor
	}
	return s.plantrepo.Delete(ctx, id)
}

func (s *PlantService) UploadPlantPhoto(ctx context.Context, id uuid.UUID, fdata models.FileData, description string) error {
	user := authservice.UserFromContext(ctx)
	if user == nil {
		return ErrNotAuthorized
	}
	if !user.HasAuthorRights() {
		return ErrNotAuthor
	}
	file, err := s.filerepo.Upload(ctx, &fdata)
	if err != nil {
		return err
	}
	_, err = s.plantrepo.Update(ctx, id, func(p *plant.Plant) (*plant.Plant, error) {
		photo, err := plant.NewPlantPhoto(file.ID, description)
		if err != nil {
			return nil, err
		}
		err = p.AddPhoto(photo)
		return p, err
	})
	return err
}

func (s *PlantService) GetPlantCategory(ctx context.Context, name string) (*plant.PlantCategory, error) {
	user := authservice.UserFromContext(ctx)
	if user == nil {
		return nil, ErrNotAuthorized
	}
	if !user.HasAuthorRights() {
		return nil, ErrNotAuthor
	}
	return s.categoryrepo.GetCategory(ctx, name)
}

func (s *PlantService) ListCategories(ctx context.Context) ([]plant.PlantCategory, error) {
	user := authservice.UserFromContext(ctx)
	if user == nil {
		return nil, ErrNotAuthorized
	}
	if !user.HasAuthorRights() {
		return nil, ErrNotAuthor
	}
	return s.categoryrepo.GetCategories(ctx)
}
