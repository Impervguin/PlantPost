package filters

import (
	_ "PlantSite/internal/infra/filters/plant-filters"
	_ "PlantSite/internal/infra/filters/post-filters"
	registry "PlantSite/internal/infra/filters/registry"
	"PlantSite/internal/models/search"
)

type PostgresPlantSearch registry.PostgresPlantSearch

type PostgresPostSearch registry.PostgresPostSearch

var (
	NewPostgresPlantSearch = registry.NewPostgresPlantSearch
	NewPostgresPostSearch  = registry.NewPostgresPostSearch
)

func MapPlantFilter(filter search.PlantFilter) (registry.PostgresPlantFilter, error) {
	return registry.MapPlantFilter(filter)
}

func MapPostFilter(filter search.PostFilter) (registry.PostgresPostFilter, error) {
	return registry.MapPostFilter(filter)
}
