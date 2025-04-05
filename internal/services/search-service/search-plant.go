package searchservice

import (
	"PlantSite/internal/models"
	"PlantSite/internal/models/plant"
	"PlantSite/internal/models/search"
	"context"
	"time"

	"github.com/google/uuid"
)

type SearchPlant struct {
	ID            uuid.UUID
	Name          string
	LatinName     string
	Description   string
	MainPhoto     models.File
	Category      string
	Specification plant.PlantSpecification
	CreatedAt     time.Time
}

func (s *SearchService) SearchPlants(ctx context.Context, plSearch *search.PlantSearch) ([]*SearchPlant, error) {
	plants, err := s.searchRepo.SearchPlants(ctx, plSearch)
	if err != nil {
		return nil, Wrap(err)
	}
	searchPlants := make([]*SearchPlant, 0, len(plants))
	for _, p := range plants {
		mainPhoto, err := s.plantFileRepo.Get(ctx, p.MainPhotoID())
		if err != nil {
			return nil, Wrap(err)
		}
		searchPlants = append(searchPlants, &SearchPlant{
			ID:            p.ID(),
			Name:          p.GetName(),
			LatinName:     p.GetLatinName(),
			Description:   p.GetDescription(),
			MainPhoto:     *mainPhoto,
			Category:      p.GetCategory(),
			Specification: p.GetSpecification(),
			CreatedAt:     p.CreatedAt(),
		})
	}
	return searchPlants, nil
}
