package plantstorage

import "PlantSite/internal/models/plant"

var globalMapper PlantSpecificationMapper

func init() {
	globalMapper = RWMapPlantSpecificationMapperFactory()
}

func Register(category string, toDb FromDb, toDomain FromDomain) {
	globalMapper.Register(category, toDb, toDomain)
}

func SpecificationFromDB(category string, json JsonB) (PlantSpecification, error) {
	return globalMapper.FromDB(category, json)
}

func SpecificationFromDomain(category string, spec plant.PlantSpecification) (PlantSpecification, error) {
	return globalMapper.FromDomain(category, spec)
}
