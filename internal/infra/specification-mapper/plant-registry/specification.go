package registry

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
type FromDB func(JsonB) (PlantSpecification, error)
