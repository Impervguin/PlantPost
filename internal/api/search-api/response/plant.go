package response

import plantspec "PlantSite/internal/api/plant-spec"

type SearchPlantItem struct {
	ID            string                       `json:"id" form:"id" binding:"required"`
	Name          string                       `json:"name" form:"name" binding:"required"`
	LatinName     string                       `json:"latin_name" form:"latin_name" binding:"required"`
	Description   string                       `json:"description" form:"description" binding:"required"`
	MainPhotoKey  string                       `json:"main_photo_key" form:"main_photo_key" binding:"required"`
	Category      string                       `json:"category" form:"category" binding:"required"`
	Specification plantspec.PlantSpecification `json:"specification" form:"specification" binding:"required"`
	CreatedAt     string                       `json:"created_at" form:"created_at" binding:"required"`
}

type SearchPlantResponse []SearchPlantItem
