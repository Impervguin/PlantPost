package postservice

import (
	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/post"
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
	if !user.HasAuthorRights() {
		return nil, auth.ErrNoAuthorRights
	}

	p, err := s.postRepo.Update(ctx, id, func(p *post.Post) (*post.Post, error) {
		if user.ID() != p.AuthorID() {
			return nil, ErrNotAuthor
		}
		err := p.UpdateContent(data.Content)
		if err != nil {
			return nil, err
		}
		err = p.UpdateTitle(data.Title)
		if err != nil {
			return nil, err
		}
		err = p.UpdateTags(data.Tags)
		if err != nil {
			return nil, err
		}
		return p, nil
	})

	return p, err
}
