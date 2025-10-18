package postservice

import (
	"PlantSite/internal/models"
	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/post"
	"PlantSite/internal/utils/logs"
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
	logs.Debugf("PostService.CreatePost: got user=%s", user.ID())
	if !user.HasAuthorRights() {
		return nil, auth.ErrNoAuthorRights
	}
	logs.Debugf("PostService.CreatePost: user %s has author rights", user.ID())

	photos := post.NewPostPhotos()
	for i, file := range files {
		if file.ContentType != "image/jpeg" && file.ContentType != "image/png" {
			return nil, Wrap(ErrInvalidFileContentType)
		}
		logs.Debugf("PostService.CreatePost: uploading photo %d for post %s", i, data.Title)
		f, err := s.fileRepo.Upload(ctx, &file)
		if err != nil {
			return nil, Wrap(err)
		}
		logs.Debugf("PostService.CreatePost: uploaded photo %s for post %s", f.ID, data.Title)
		photo, err := post.NewPostPhoto(f.ID, i+1)
		if err != nil {
			return nil, Wrap(err)
		}
		logs.Debugf("PostService.CreatePost: created photo %s for post %s", f.ID, data.Title)
		// photos = append(photos, *photo)
		err = photos.Add(photo)
		if err != nil {
			return nil, Wrap(err)
		}
		logs.Debugf("PostService.CreatePost: added photo %s for post %s", f.ID, data.Title)
	}

	post, err := post.NewPost(data.Title, data.Content, data.Tags, user.ID(), photos)
	if err != nil {
		return nil, err
	}
	logs.Debugf("PostService.CreatePost: created post %s", post.ID())

	return s.postRepo.Create(ctx, post)

}
