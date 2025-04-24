package request

import (
	plantfilters "PlantSite/internal/api/search-api/plant-filters"
	postfilters "PlantSite/internal/api/search-api/post-filters"

	"github.com/google/uuid"
)

type SearchPostsItem postfilters.PostFilter

type SearchPostsRequest []SearchPostsItem

type SearchPlantsItem plantfilters.PlantFilter

type SearchPlantsRequest []SearchPlantsItem

type GetPlantRequest struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

type GetPostRequest struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}
