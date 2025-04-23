package request

import (
	"PlantSite/internal/api/plant-api/spec"

	"github.com/google/uuid"
)

type CreatePlantRequest struct {
	Name        string
	LatinName   string
	Description string
	Category    string
	Spec        spec.PlantSpecification
}

type GetPlantRequest struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

type UpdatePlantSpecRequest struct {
	ID       uuid.UUID               `uri:"id" binding:"required"`
	Category string                  `json:"category" form:"category" binding:"required"`
	Spec     spec.PlantSpecification `json:"specification" form:"specification" binding:"required"`
}

type DeletePlantRequest struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

type UploadPlantPhotoRequest struct {
	ID          uuid.UUID `uri:"id" binding:"required"`
	Description string    `json:"description" form:"description" binding:"required"`
}
