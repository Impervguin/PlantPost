package request

import (
	plantfilters "PlantSite/internal/api/search-api/plant-filters"
	postfilters "PlantSite/internal/api/search-api/post-filters"
)

type SearchPostsItem postfilters.PostFilter

type SearchPostsRequest []SearchPostsItem

type SearchPlantsItem plantfilters.PlantFilter

type SearchPlantsRequest []SearchPlantsItem
