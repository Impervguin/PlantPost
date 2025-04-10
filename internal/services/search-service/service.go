package searchservice

import (
	"PlantSite/internal/models"
	"PlantSite/internal/models/search"
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
