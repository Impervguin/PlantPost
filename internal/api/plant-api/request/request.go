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
	Category string                  `         binding:"required" json:"category"      form:"category"`
	Spec     spec.PlantSpecification `         binding:"required" json:"specification" form:"specification"`
}

type DeletePlantRequest struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

type UploadPlantPhotoRequest struct {
	ID          uuid.UUID `uri:"id" binding:"required"`
	Description string    `         binding:"required" json:"description" form:"description"`
}
