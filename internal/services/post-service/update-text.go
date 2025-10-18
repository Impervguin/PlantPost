package postservice

import (
	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/post"
	"PlantSite/internal/utils/logs"
	"context"

	"github.com/google/uuid"
)

type UpdatePostTextData struct {
	Title   string
	Content post.Content
	Tags    []string
}

func (s *PostService) UpdatePost(ctx context.Context, id uuid.UUID, data UpdatePostTextData) (*post.Post, error) {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return nil, auth.ErrNotAuthorized
	}
	logs.Debugf("PostService.UpdatePost: got user=%s", user.ID())
	if !user.HasAuthorRights() {
		return nil, auth.ErrNoAuthorRights
	}
	logs.Debugf("PostService.UpdatePost: user %s has author rights", user.ID())

	p, err := s.postRepo.Update(ctx, id, func(p *post.Post) (*post.Post, error) {
		if user.ID() != p.AuthorID() {
			return nil, ErrNotAuthor
		}
		logs.Debugf("PostService.UpdatePost: %s author for post %s", user.ID(), p.ID())
		err := p.UpdateContent(data.Content)
		if err != nil {
			return nil, err
		}
		logs.Debugf("PostService.UpdatePost: updated content for post %s", p.ID())
		err = p.UpdateTitle(data.Title)
		if err != nil {
			return nil, err
		}
		logs.Debugf("PostService.UpdatePost: updated title for post %s", p.ID())
		err = p.UpdateTags(data.Tags)
		if err != nil {
			return nil, err
		}
		logs.Debugf("PostService.UpdatePost: updated tags for post %s", p.ID())
		return p, nil
	})

	return p, err
}
