package mapper

import (
	plantfilters "PlantSite/internal/api/search-api/plant-filters"
	postfilters "PlantSite/internal/api/search-api/post-filters"
	"PlantSite/internal/api/search-api/request"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

type GetPlantRequest struct {
	ID string `uri:"id" binding:"required"`
}

func MapGetPlantRequest(c *gin.Context) (*request.GetPlantRequest, error) {
	var req GetPlantRequest
	if err := c.ShouldBindUri(&req); err != nil {
		return nil, fmt.Errorf("can't bind uri: %w", err)
	}
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, fmt.Errorf("can't parse id: %w", err)
	}
	return &request.GetPlantRequest{
		ID: id,
	}, nil
}

type GetPostRequest struct {
	ID string `uri:"id" binding:"required"`
}

func MapGetPostRequest(c *gin.Context) (*request.GetPostRequest, error) {
	var req GetPostRequest
	if err := c.ShouldBindUri(&req); err != nil {
		return nil, fmt.Errorf("can't bind uri: %w", err)
	}
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, fmt.Errorf("can't parse id: %w", err)
	}
	return &request.GetPostRequest{
		ID: id,
	}, nil
}
