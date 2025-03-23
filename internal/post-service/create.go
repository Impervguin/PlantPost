package postservice

import (
	authservice "PlantSite/internal/auth-service"
	"PlantSite/internal/models"
	"PlantSite/internal/models/post"
	"context"
)

type CreatePostTextData struct {
	Title   string
	Content post.Content
	Tags    []string
}

func (s *PostService) CreatePost(ctx context.Context, data CreatePostTextData, files []models.FileData) (*post.Post, error) {
	user := authservice.UserFromContext(ctx)
	if user == nil {
		return nil, ErrNotAuthorized
	}
	if !user.HasAuthorRights() {
		return nil, ErrNotAuthor
	}

	// photos := make([]post.PostPhoto, 0, len(files))
	photos := post.NewPostPhotos()
	for i, file := range files {
		if file.ContentType != "image/jpeg" && file.ContentType != "image/png" {
			return nil, Wrap(ErrInvalidFileContentType)
		}
		f, err := s.fileRepo.Upload(ctx, &file)
		if err != nil {
			return nil, err
		}
		photo, err := post.NewPostPhoto(f.ID, i+1)
		if err != nil {
			return nil, err
		}
		// photos = append(photos, *photo)
		photos.Add(photo)
	}

	post, err := post.NewPost(data.Title, data.Content, data.Tags, user.ID(), photos)
	if err != nil {
		return nil, err
	}

	return s.postRepo.Create(ctx, post)

}
