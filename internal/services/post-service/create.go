package postservice

import (
	"PlantSite/internal/models"
	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/post"
	"context"
)

type CreatePostTextData struct {
	Title   string
	Content post.Content
	Tags    []string
}

func (s *PostService) CreatePost(ctx context.Context, data CreatePostTextData, files []models.FileData) (*post.Post, error) {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return nil, auth.ErrNotAuthorized
	}
	if !user.HasAuthorRights() {
		return nil, auth.ErrNoAuthorRights
	}

	// photos := make([]post.PostPhoto, 0, len(files))
	photos := post.NewPostPhotos()
	for i, file := range files {
		if file.ContentType != "image/jpeg" && file.ContentType != "image/png" {
			return nil, Wrap(ErrInvalidFileContentType)
		}
		f, err := s.fileRepo.Upload(ctx, &file)
		if err != nil {
			return nil, Wrap(err)
		}
		photo, err := post.NewPostPhoto(f.ID, i+1)
		if err != nil {
			return nil, Wrap(err)
		}
		// photos = append(photos, *photo)
		err = photos.Add(photo)
		if err != nil {
			return nil, Wrap(err)
		}
	}

	post, err := post.NewPost(data.Title, data.Content, data.Tags, user.ID(), photos)
	if err != nil {
		return nil, err
	}

	return s.postRepo.Create(ctx, post)

}
