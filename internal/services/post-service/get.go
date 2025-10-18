package postservice

import (
	"PlantSite/internal/models"
	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/post"
	"PlantSite/internal/utils/logs"
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
		return nil, auth.ErrNotAuthorized
	}
	logs.Debugf("PostService.GetPost: got user=%s", user.ID())
	if !user.HasAuthorRights() {
		return nil, auth.ErrNoAuthorRights
	}
	logs.Debugf("PostService.GetPost: user %s has author rights", user.ID())
	if id == uuid.Nil {
		return nil, fmt.Errorf("id must be non-nil")
	}
	post, err := s.postRepo.Get(ctx, id)
	if err != nil {
		return nil, Wrap(err)
	}
	logs.Debugf("PostService.GetPost: got post %s", post.ID())
	photos := make([]GetPostPhoto, 0)

	for _, p := range post.Photos().List() {
		file, err := s.fileRepo.Get(ctx, p.FileID())
		if err != nil {
			return nil, Wrap(err)
		}
		logs.Debugf("PostService.GetPost: got photo %s for post %s", file.ID, post.ID())
		photos = append(photos, GetPostPhoto{
			ID:          p.ID(),
			PlaceNumber: p.PlaceNumber(),
			File:        *file,
		})
		logs.Debugf("PostService.GetPost: got photo %s for post %s", file.ID, post.ID())
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
