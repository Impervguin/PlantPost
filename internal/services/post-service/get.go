package postservice

import (
	"PlantSite/internal/models"
	"PlantSite/internal/models/post"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type GetPost struct {
	ID      uuid.UUID
	Title   string
	Content post.Content
	Tags    []string
	Photos  []GetPostPhoto

	AuthorID  uuid.UUID
	UpdatedAt time.Time
	CreatedAt time.Time
}

type GetPostPhoto struct {
	ID          uuid.UUID
	PlaceNumber int
	File        models.File
}

func (s *PostService) GetPost(ctx context.Context, id uuid.UUID) (*GetPost, error) {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return nil, ErrNotAuthorized
	}
	if !user.HasAuthorRights() {
		return nil, ErrNotAuthor
	}
	if id == uuid.Nil {
		return nil, fmt.Errorf("id must be non-nil")
	}
	post, err := s.postRepo.Get(ctx, id)
	if err != nil {
		return nil, Wrap(err)
	}
	photos := make([]GetPostPhoto, 0)

	for _, p := range post.Photos().List() {
		file, err := s.fileRepo.Get(ctx, p.FileID())
		if err != nil {
			return nil, Wrap(err)
		}
		photos = append(photos, GetPostPhoto{
			ID:          p.ID(),
			PlaceNumber: p.PlaceNumber(),
			File:        *file,
		})
	}

	return &GetPost{
		ID:      post.ID(),
		Title:   post.Title(),
		Content: post.Content(),
		Tags:    post.Tags(),
		Photos:  photos,

		AuthorID:  post.AuthorID(),
		UpdatedAt: post.UpdatedAt(),
		CreatedAt: post.CreatedAt(),
	}, nil
}
