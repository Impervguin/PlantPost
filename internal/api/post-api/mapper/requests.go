package mapper

import (
	"PlantSite/internal/api/post-api/request"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetPostRequest struct {
	ID string `uri:"id" binding:"required"`
}

func MapPostGetRequest(c *gin.Context) (*request.GetPostRequest, error) {
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

type DeletePostRequest struct {
	ID string `uri:"id" binding:"required"`
}

func MapPostDeleteRequest(c *gin.Context) (*request.DeletePostRequest, error) {
	var req DeletePostRequest
	if err := c.ShouldBindUri(&req); err != nil {
		return nil, fmt.Errorf("can't bind uri: %w", err)
	}
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, fmt.Errorf("can't parse id: %w", err)
	}
	return &request.DeletePostRequest{
		ID: id,
	}, nil
}

type UpdatePostRequestID struct {
	ID string `uri:"id" binding:"required"`
}

type UpdatePostRequestBody struct {
	Title   string   `json:"title" form:"title" binding:"required"`
	Content string   `json:"content" form:"content" binding:"required"`
	Tags    []string `json:"tags" form:"tags"`
}

func MapPostUpdateRequest(c *gin.Context) (*request.UpdateTextPostRequest, error) {
	var reqID UpdatePostRequestID
	if err := c.ShouldBindUri(&reqID); err != nil {
		return nil, fmt.Errorf("can't bind uri: %w", err)
	}
	id, err := uuid.Parse(reqID.ID)
	if err != nil {
		return nil, fmt.Errorf("can't parse id: %w", err)
	}
	var req UpdatePostRequestBody
	if err := c.ShouldBind(&req); err != nil {
		return nil, fmt.Errorf("can't bind body: %w", err)
	}

	if req.Tags == nil {
		req.Tags = make([]string, 0)
	}
	return &request.UpdateTextPostRequest{
		ID:      id,
		Title:   req.Title,
		Content: req.Content,
		Tags:    req.Tags,
	}, nil
}
