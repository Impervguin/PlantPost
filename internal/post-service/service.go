package postservice

import (
	"PlantSite/internal/models"
	"PlantSite/internal/models/post"
)

type PostService struct {
	postRepo post.PostRepository
	fileRepo models.FileRepository
}

func NewPostService(repo post.PostRepository, fileRepo models.FileRepository) *PostService {
	return &PostService{postRepo: repo, fileRepo: fileRepo}
}


