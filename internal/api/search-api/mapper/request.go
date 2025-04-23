package mapper

import (
	plantfilters "PlantSite/internal/api/search-api/plant-filters"
	postfilters "PlantSite/internal/api/search-api/post-filters"
	"PlantSite/internal/api/search-api/request"

	"github.com/gin-gonic/gin"
)

type SearchPostsItem struct {
	Type string                 `json:"type" form:"type" binding:"required"`
	Data map[string]interface{} `json:"params" form:"params" binding:"required"`
}

func MapSearchPostsRequest(c *gin.Context) (request.SearchPostsRequest, error) {
	var items []SearchPostsItem
	if err := c.ShouldBind(&items); err != nil {
		return nil, err
	}
	req := make(request.SearchPostsRequest, 0, len(items))
	for _, item := range items {
		f, err := postfilters.ParsePostFilter(item.Type, item.Data)
		if err != nil {
			return nil, err
		}
		req = append(req, f)
	}
	return req, nil
}

type SearchPlantsItem struct {
	Type string                 `json:"type" form:"type" binding:"required"`
	Data map[string]interface{} `json:"params" form:"params" binding:"required"`
}

func MapSearchPlantsRequest(c *gin.Context) (request.SearchPlantsRequest, error) {
	var items []SearchPlantsItem
	if err := c.ShouldBind(&items); err != nil {
		return nil, err
	}
	req := make(request.SearchPlantsRequest, 0, len(items))
	for _, item := range items {
		f, err := plantfilters.ParsePlantFilter(item.Type, item.Data)
		if err != nil {
			return nil, err
		}
		req = append(req, f)
	}
	return req, nil
}
