package response

import "PlantSite/internal/api/plant-api/spec"

type GetPlantResponse struct {
	ID            string `json:"id" form:"id" binding:"required"`
	Name          string `json:"name" form:"name" binding:"required"`
	LatinName     string `json:"latin_name" form:"latin_name" binding:"required"`
	Description   string `json:"description" form:"description" binding:"required"`
	MainPhotoKey  string `json:"main_photo_key" form:"main_photo_key" binding:"required"`
	Photos        []GetPlantPhoto
	Category      string `json:"category" form:"category" binding:"required"`
	Specification spec.PlantSpecification
	CreatedAt     string `json:"created_at" form:"created_at" binding:"required"`
}

type GetPlantPhoto struct {
	ID          string `json:"id" form:"id" binding:"required"`
	Key         string `json:"key" form:"key" binding:"required"`
	Description string `json:"description" form:"description" binding:"required"`
}
