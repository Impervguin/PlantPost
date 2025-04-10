package specificationmapper

import (
	registry "PlantSite/internal/infra/specification-mapper/plant-registry"
	_ "PlantSite/internal/infra/specification-mapper/specification"
	"PlantSite/internal/models/plant"
)

type JsonB registry.JsonB

type PlantSpecification registry.PlantSpecification

func SpecificationFromDB(category string, json JsonB) (PlantSpecification, error) {
	return registry.MapFromDB(category, registry.JsonB(json))
}

func SpecificationFromDomain(category string, spec plant.PlantSpecification) (PlantSpecification, error) {
	return registry.MapFromDomain(category, spec)
}
