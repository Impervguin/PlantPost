package mapper

import (
	"PlantSite/internal/api/album-api/request"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateAlbumRequest struct {
	Name        string   `json:"name" form:"name" binding:"required"`
	Description string   `json:"description" form:"description" binding:"required"`
	PlantIDs    []string `json:"plant_ids" form:"plant_ids" binding:"required"`
}

func MapCreateAlbumRequest(c *gin.Context) (*request.CreateAlbumRequest, error) {
	var req CreateAlbumRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}
	plantIDs := make(uuid.UUIDs, 0, len(req.PlantIDs))
	for _, id := range req.PlantIDs {
		pid, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		plantIDs = append(plantIDs, pid)
	}
	return &request.CreateAlbumRequest{
		Name:        req.Name,
		Description: req.Description,
		PlantIDs:    plantIDs,
	}, nil
}

type GetAlbumRequest struct {
	ID string `uri:"id" binding:"required"`
}

func MapGetAlbumRequest(c *gin.Context) (*request.GetAlbumRequest, error) {
	var req GetAlbumRequest
	if err := c.ShouldBindUri(&req); err != nil {
		return nil, fmt.Errorf("can't bind uri: %w", err)
	}
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, fmt.Errorf("can't parse id: %w", err)
	}
	return &request.GetAlbumRequest{
		ID: id,
	}, nil
}

type UpdateAlbumRequestID struct {
	ID string `uri:"id" binding:"required"`
}

func fetchAlbumID(c *gin.Context) (*UpdateAlbumRequestID, error) {
	var req UpdateAlbumRequestID
	if err := c.ShouldBindUri(&req); err != nil {
		return nil, fmt.Errorf("can't bind uri: %w", err)
	}
	return &req, nil
}

type UpdateAlbumNameRequest struct {
	Name string `json:"name" form:"name" binding:"required"`
}

func MapUpdateAlbumNameRequest(c *gin.Context) (*request.UpdateAlbumNameRequest, error) {
	req, err := fetchAlbumID(c)
	if err != nil {
		return nil, err
	}
	var reqName UpdateAlbumNameRequest
	if err := c.ShouldBind(&reqName); err != nil {
		return nil, fmt.Errorf("can't bind name: %w", err)
	}
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, fmt.Errorf("can't parse id: %w", err)
	}

	return &request.UpdateAlbumNameRequest{
		ID:   id,
		Name: reqName.Name,
	}, nil
}

type UpdateAlbumDescriptionRequest struct {
	Description string `json:"description" form:"description" binding:"required"`
}

func MapUpdateAlbumDescriptionRequest(c *gin.Context) (*request.UpdateAlbumDescriptionRequest, error) {
	req, err := fetchAlbumID(c)
	if err != nil {
		return nil, err
	}
	var reqDesc UpdateAlbumDescriptionRequest
	if err := c.ShouldBind(&reqDesc); err != nil {
		return nil, fmt.Errorf("can't bind description: %w", err)
	}
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, fmt.Errorf("can't parse id: %w", err)
	}
	return &request.UpdateAlbumDescriptionRequest{
		ID:          id,
		Description: reqDesc.Description,
	}, nil
}

type AddPlantToAlbumRequest struct {
	PlantID string `json:"plant_id" form:"plant_id" binding:"required"`
}

func MapAddPlantToAlbumRequest(c *gin.Context) (*request.AddPlantToAlbumRequest, error) {
	var reqID UpdateAlbumRequestID
	if err := c.ShouldBindUri(&reqID); err != nil {
		return nil, fmt.Errorf("can't bind uri: %w", err)
	}
	var req AddPlantToAlbumRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, fmt.Errorf("can't bind body: %w", err)
	}
	id, err := uuid.Parse(reqID.ID)
	if err != nil {
		return nil, fmt.Errorf("can't parse id: %w", err)
	}
	plantID, err := uuid.Parse(req.PlantID)
	if err != nil {
		return nil, fmt.Errorf("can't parse plant id: %w", err)
	}
	return &request.AddPlantToAlbumRequest{
		ID:      id,
		PlantID: plantID,
	}, nil
}

type RemovePlantFromAlbumRequest struct {
	PlantID string `json:"plant_id" form:"plant_id" binding:"required"`
}

func MapRemovePlantFromAlbumRequest(c *gin.Context) (*request.RemovePlantFromAlbumRequest, error) {
	var reqID UpdateAlbumRequestID
	if err := c.ShouldBindUri(&reqID); err != nil {
		return nil, fmt.Errorf("can't bind uri: %w", err)
	}
	var req RemovePlantFromAlbumRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, fmt.Errorf("can't bind body: %w", err)
	}
	id, err := uuid.Parse(reqID.ID)
	if err != nil {
		return nil, fmt.Errorf("can't parse id: %w", err)
	}
	plantID, err := uuid.Parse(req.PlantID)
	if err != nil {
		return nil, fmt.Errorf("can't parse plant id: %w", err)
	}
	return &request.RemovePlantFromAlbumRequest{
		ID:      id,
		PlantID: plantID,
	}, nil
}

type DeleteAlbumRequest struct {
	ID string `uri:"id" binding:"required"`
}

func MapDeleteAlbumRequest(c *gin.Context) (*request.DeleteAlbumRequest, error) {
	var req DeleteAlbumRequest
	if err := c.ShouldBindUri(&req); err != nil {
		return nil, fmt.Errorf("can't bind uri: %w", err)
	}
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, fmt.Errorf("can't parse id: %w", err)
	}
	return &request.DeleteAlbumRequest{
		ID: id,
	}, nil
}
