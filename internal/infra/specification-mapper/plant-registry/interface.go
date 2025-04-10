package registry

import (
	"PlantSite/internal/models/plant"
)

var (
	globalRegistry *PlantSpecificationRegistry = NewPlantSpecificationRegistry()
)

func Register(category string, fromDB FromDB, fromDomain FromDomain) error {
	return globalRegistry.Register(category, fromDB, fromDomain)
}

func MapFromDB(category string, json JsonB) (PlantSpecification, error) {
	return globalRegistry.FromDB(category, json)
}

func MapFromDomain(category string, spec plant.PlantSpecification) (PlantSpecification, error) {
	return globalRegistry.FromDomain(category, spec)
}
