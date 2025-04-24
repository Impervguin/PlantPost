package response

import (
	"github.com/google/uuid"
)

type GetPostResponse struct {
	ID          uuid.UUID      `json:"id"`
	Title       string         `json:"title"`
	Content     string         `json:"content"`
	ContentType string         `json:"content_type"`
	Tags        []string       `json:"tags"`
	Photos      []GetPostPhoto `json:"photos"`

	AuthorID  uuid.UUID `json:"author_id"`
	UpdatedAt string    `json:"updated_at"`
	CreatedAt string    `json:"created_at"`
}

type GetPostPhoto struct {
	ID          uuid.UUID `json:"id"`
	PlaceNumber int       `json:"place_number"`
	Key         string    `json:"key"`
}
