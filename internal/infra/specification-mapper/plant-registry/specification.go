package registry

import (
	"PlantSite/internal/models/plant"
	"errors"
)

var (
	ErrCategoryNotFound = errors.New("category not found in mapper")
)

type PlantSpecification interface {
	ToJSONB() (JSONB, error)
	ToDomain() (plant.PlantSpecification, error)
}

type FromDomain func(plant.PlantSpecification) (PlantSpecification, error)
type FromDB func(JSONB) (PlantSpecification, error)
