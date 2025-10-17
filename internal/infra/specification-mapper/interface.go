package specificationmapper

import (
	registry "PlantSite/internal/infra/specification-mapper/plant-registry"
	_ "PlantSite/internal/infra/specification-mapper/specification"
	"PlantSite/internal/models/plant"
)

type JSONB registry.JSONB

type PlantSpecification registry.PlantSpecification

func SpecificationFromDB(category string, json JSONB) (PlantSpecification, error) {
	return registry.MapFromDB(category, registry.JSONB(json))
}

func SpecificationFromDomain(category string, spec plant.PlantSpecification) (PlantSpecification, error) {
	return registry.MapFromDomain(category, spec)
}
