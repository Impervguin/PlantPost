package searchservice

import (
	"PlantSite/internal/models"
	"PlantSite/internal/models/post"
	"PlantSite/internal/utils/logs"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s *SearchService) GetPostByID(ctx context.Context, id uuid.UUID) (*post.Post, error) {
	return s.searchRepo.GetPostByID(ctx, id)
}

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

func (s *SearchService) GetPost(ctx context.Context, id uuid.UUID) (*GetPost, error) {
	if id == uuid.Nil {
		return nil, Wrap(fmt.Errorf("id must be non-nil"))
	}
	logs.Debugf("SearchService.GetPost: got post ID %s", id)
	post, err := s.searchRepo.GetPostByID(ctx, id)
	if err != nil {
		return nil, Wrap(err)
	}
	logs.Debugf("SearchService.GetPost: got post %s", post.ID())
	photos := make([]GetPostPhoto, 0)

	for _, p := range post.Photos().List() {
		file, err := s.postFileRepo.Get(ctx, p.FileID())
		if err != nil {
			return nil, Wrap(err)
		}
		photos = append(photos, GetPostPhoto{
			ID:          p.ID(),
			PlaceNumber: p.PlaceNumber(),
			File:        *file,
		})
	}
	logs.Debugf("SearchService.GetPost: got %d photos", len(photos))

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
