package request

import "github.com/google/uuid"

type CreatePostRequest struct {
	Title   string   `json:"title" form:"title" binding:"required"`
	Content string   `json:"content" form:"content" binding:"required"`
	Tags    []string `json:"tags" form:"tags"`
}

type GetPostRequest struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

type DeletePostRequest struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

type UpdateTextPostRequest struct {
	ID      uuid.UUID `uri:"id" binding:"required"`
	Title   string    `json:"title" form:"title" binding:"required"`
	Content string    `json:"content" form:"content" binding:"required"`
	Tags    []string  `json:"tags" form:"tags"`
}
