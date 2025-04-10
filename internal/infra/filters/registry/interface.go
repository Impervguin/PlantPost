package filters

import (
	"PlantSite/internal/models/search"
)

var (
	globalRegistry *FilterRegistry = NewFilterRegistry()
)

func RegisterPlantFilter(id string, factory PlantFilterFactory) error {
	return globalRegistry.RegisterPlantFilter(id, factory)
}

func RegisterPostFilter(id string, factory PostFilterFactory) error {
	return globalRegistry.RegisterPostFilter(id, factory)
}

func MapPlantFilter(filter search.PlantFilter) (PostgresPlantFilter, error) {
	return globalRegistry.MapPlantFilter(filter)
}

func MapPostFilter(filter search.PostFilter) (PostgresPostFilter, error) {
	return globalRegistry.MapPostFilter(filter)
}
