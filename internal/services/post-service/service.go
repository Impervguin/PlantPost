package postservice

import (
	"PlantSite/internal/models"
	"PlantSite/internal/models/post"
	authservice "PlantSite/internal/services/auth-service"
)

type PostService struct {
	postRepo post.PostRepository
	fileRepo models.FileRepository
	auth     *authservice.AuthService
}

func NewPostService(repo post.PostRepository, fileRepo models.FileRepository, auth *authservice.AuthService) *PostService {
	if repo == nil {
		panic("nil repository")
	}
	if fileRepo == nil {
		panic("nil file repository")
	}
	if auth == nil {
		panic("nil auth")
	}
	return &PostService{postRepo: repo, fileRepo: fileRepo, auth: auth}
}
