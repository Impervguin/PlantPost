package mapper

import (
	"PlantSite/internal/api/plant-api/request"
	"PlantSite/internal/api/plant-api/spec"
	"PlantSite/internal/models/plant"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreatePlantBase struct {
	Name        string `json:"name" form:"name" binding:"required"`
	LatinName   string `json:"latin_name" form:"latin_name" binding:"required"`
	Description string `json:"description" form:"description" binding:"required"`
	Category    string `json:"category" form:"category" binding:"required"`
}

type CreateConiferousPlantRequest struct {
	Spec spec.ConiferousSpecification `json:"specification" form:"specification" binding:"required"`
}

type CreateDeciduousPlantRequest struct {
	Spec spec.DeciduousSpecification `json:"specification" form:"specification" binding:"required"`
}

type PlantCategory struct {
	Category string `json:"category" form:"category" binding:"required"`
}

func MapCreatePlantRequest(c *gin.Context) (*request.CreatePlantRequest, error) {
	var reqBase CreatePlantBase
	if err := c.ShouldBind(&reqBase); err != nil {
		return nil, fmt.Errorf("can't bind base: %w", err)
	}
	var reqSpec spec.PlantSpecification
	switch reqBase.Category {
	case plant.ConiferousCategory:
		var req CreateConiferousPlantRequest
		if err := c.ShouldBind(&req); err != nil {
			return nil, fmt.Errorf("can't bind request: %w", err)
		}
		reqSpec = &req.Spec
	case plant.DeciduousCategory:
		var req CreateDeciduousPlantRequest
		if err := c.ShouldBind(&req); err != nil {
			return nil, fmt.Errorf("can't bind request: %w", err)
		}
		reqSpec = &req.Spec
	default:
		return nil, fmt.Errorf("invalid category: %s", reqBase.Category)
	}
	return &request.CreatePlantRequest{
		Name:        reqBase.Name,
		LatinName:   reqBase.LatinName,
		Description: reqBase.Description,
		Category:    reqBase.Category,
		Spec:        reqSpec,
	}, nil
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

type UpdatespecRequestID struct {
	ID string `uri:"id" binding:"required"`
}

type UpdatespecRequestCategory struct {
	Category string `json:"category" form:"category" binding:"required"`
}

type UpdatePlantConiferousRequestBody struct {
	Spec spec.ConiferousSpecification `json:"specification" form:"specification" binding:"required"`
}

type UpdatePlantDeciduousRequestBody struct {
	Spec spec.DeciduousSpecification `json:"specification" form:"specification" binding:"required"`
}

func MapUpdatePlantSpecRequest(c *gin.Context) (*request.UpdatePlantSpecRequest, error) {
	var reqID UpdatespecRequestID
	if err := c.ShouldBindUri(&reqID); err != nil {
		return nil, fmt.Errorf("can't bind uri: %w", err)
	}
	id, err := uuid.Parse(reqID.ID)
	if err != nil {
		return nil, fmt.Errorf("can't parse id: %w", err)
	}
	var req UpdatespecRequestCategory
	if err := c.ShouldBind(&req); err != nil {
		return nil, fmt.Errorf("can't bind category: %w", err)
	}

	var reqSpec spec.PlantSpecification
	switch req.Category {
	case plant.ConiferousCategory:
		var reqBody UpdatePlantConiferousRequestBody
		if err := c.ShouldBind(&reqBody); err != nil {
			return nil, fmt.Errorf("can't bind body: %w", err)
		}
	case plant.DeciduousCategory:
		var reqBody UpdatePlantDeciduousRequestBody
		if err := c.ShouldBind(&reqBody); err != nil {
			return nil, fmt.Errorf("can't bind body: %w", err)
		}
	default:
		return nil, fmt.Errorf("invalid category: %s", req.Category)
	}
	return &request.UpdatePlantSpecRequest{
		ID:       id,
		Category: req.Category,
		Spec:     reqSpec,
	}, nil
}

type DeletePlantRequest struct {
	ID string `uri:"id" binding:"required"`
}

func MapDeletePlantRequest(c *gin.Context) (*request.DeletePlantRequest, error) {
	var req DeletePlantRequest
	if err := c.ShouldBindUri(&req); err != nil {
		return nil, fmt.Errorf("can't bind uri: %w", err)
	}
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, fmt.Errorf("can't parse id: %w", err)
	}
	return &request.DeletePlantRequest{
		ID: id,
	}, nil
}

type UploadPlantPhotoRequestID struct {
	ID string `uri:"id" binding:"required"`
}

type UploadPlantPhotoRequestDescription struct {
	Description string `json:"description" form:"description" binding:"required"`
}

func MapUploadPlantPhotoRequest(c *gin.Context) (*request.UploadPlantPhotoRequest, error) {
	var reqID UploadPlantPhotoRequestID
	if err := c.ShouldBindUri(&reqID); err != nil {
		return nil, fmt.Errorf("can't bind uri: %w", err)
	}
	id, err := uuid.Parse(reqID.ID)
	if err != nil {
		return nil, fmt.Errorf("can't parse id: %w", err)
	}
	var req UploadPlantPhotoRequestDescription
	if err := c.ShouldBind(&req); err != nil {
		return nil, fmt.Errorf("can't bind description: %w", err)
	}
	return &request.UploadPlantPhotoRequest{
		ID:          id,
		Description: req.Description,
	}, nil
}
