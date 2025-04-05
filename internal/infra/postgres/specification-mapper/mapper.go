package plantstorage

import (
	"PlantSite/internal/models/plant"
	"errors"
)

var (
	ErrCategoryNotFound = errors.New("category not found in mapper")
)

type PlantSpecification interface {
	ToJsonB() (JsonB, error)
	ToDomain() (plant.PlantSpecification, error)
}

type FromDomain func(plant.PlantSpecification) (PlantSpecification, error)
type FromDb func(JsonB) (PlantSpecification, error)

type PlantSpecificationMapperFactory func() PlantSpecificationMapper

type PlantSpecificationMapper interface {
	FromDomain(string, plant.PlantSpecification) (PlantSpecification, error)
	FromDB(string, JsonB) (PlantSpecification, error)
	Register(string, FromDb, FromDomain)
}
