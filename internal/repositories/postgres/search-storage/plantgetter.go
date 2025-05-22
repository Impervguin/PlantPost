package searchstorage

import (
	"PlantSite/internal/models/plant"
	"PlantSite/internal/models/search"
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrPlantNotFound       = plant.ErrPlantNotFound
	ErrMultiplePlantsFound = errors.New("multiple plants found with the same name")
)

type SearchPlantGetter struct {
	repo *PostgresSearchRepository
}

func NewSearchPlantGetter(repo *PostgresSearchRepository) *SearchPlantGetter {
	return &SearchPlantGetter{repo: repo}
}

func (g *SearchPlantGetter) GetPlants(uuids []uuid.UUID) ([]*plant.Plant, error) {
	srch := search.NewPlantSearch()
	srch.AddFilter(search.NewPlantIDsFilter(uuids))
	return g.repo.SearchPlants(context.Background(), srch)
}

func (g *SearchPlantGetter) GetPlantByName(name string) (*plant.Plant, error) {
	srch := search.NewPlantSearch()
	srch.AddFilter(search.NewExactPlantNameFilter(name))
	plants, err := g.repo.SearchPlants(context.Background(), srch)
	if err != nil {
		return nil, err
	}
	if len(plants) == 0 {
		return nil, plant.ErrPlantNotFound
	}
	if len(plants) > 1 {
		return nil, ErrMultiplePlantsFound
	}
	return plants[0], nil
}
