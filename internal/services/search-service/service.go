package searchservice

import (
	"PlantSite/internal/models"
	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/search"
	"context"
	"fmt"
)

type SearchService struct {
	searchRepo    search.SearchRepository
	plantFileRepo models.FileRepository
	postFileRepo  models.FileRepository
}

func NewSearchService(repo search.SearchRepository, plantFileRepo models.FileRepository, postFileRepo models.FileRepository) *SearchService {
	if repo == nil {
		panic("Search repository cannot be nil")
	}
	if plantFileRepo == nil {
		panic("Plant file repository cannot be nil")
	}
	if postFileRepo == nil {
		panic("Post file repository cannot be nil")
	}
	return &SearchService{
		searchRepo:    repo,
		plantFileRepo: plantFileRepo,
		postFileRepo:  postFileRepo,
	}
}

func (s *SearchService) PostAuthors(ctx context.Context) ([]*auth.Author, error) {
	authors, err := s.searchRepo.GetPostAuthors(ctx)
	if err != nil {
		return nil, fmt.Errorf("SearchService.PostAuthors failed: %w", err)
	}

	return authors, nil
}

func (s *SearchService) PostTags(ctx context.Context) ([]string, error) {
	tags, err := s.searchRepo.GetPostTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("SearchService.PostTags failed: %w", err)
	}

	return tags, nil
}
