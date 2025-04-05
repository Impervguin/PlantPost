package searchservice

import (
	"PlantSite/internal/models"
	"PlantSite/internal/models/post"
	"PlantSite/internal/models/search"
	"context"
	"time"

	"github.com/google/uuid"
)

type SearchPost struct {
	ID      uuid.UUID
	Title   string
	Content post.Content
	Tags    []string
	Photos  []SearchPostPhoto

	AuthorID  uuid.UUID
	UpdatedAt time.Time
	CreatedAt time.Time
}

type SearchPostPhoto struct {
	ID          uuid.UUID
	PlaceNumber int
	File        models.File
}

func (s *SearchService) SearchPosts(ctx context.Context, plSearch *search.PostSearch) ([]*SearchPost, error) {
	posts, err := s.searchRepo.SearchPosts(ctx, plSearch)
	if err != nil {
		return nil, Wrap(err)
	}
	searchPosts := make([]*SearchPost, 0, len(posts))
	for _, p := range posts {
		photos := make([]SearchPostPhoto, 0)
		for _, photo := range p.Photos().List() {
			file, err := s.postFileRepo.Get(ctx, photo.FileID())
			if err != nil {
				return nil, Wrap(err)
			}
			photos = append(photos, SearchPostPhoto{
				ID:          photo.ID(),
				PlaceNumber: photo.PlaceNumber(),
				File:        *file,
			})
		}
		searchPosts = append(searchPosts, &SearchPost{
			ID:      p.ID(),
			Title:   p.Title(),
			Content: p.Content(),
			Tags:    p.Tags(),
			Photos:  photos,

			AuthorID:  p.AuthorID(),
			UpdatedAt: p.UpdatedAt(),
			CreatedAt: p.CreatedAt(),
		})
	}
	return searchPosts, nil
}
