package postservice

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *PostService) Delete(ctx context.Context, id uuid.UUID) error {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return ErrNotAuthorized
	}
	if !user.HasAuthorRights() {
		return ErrNotAuthor
	}
	if id == uuid.Nil {
		return fmt.Errorf("nil post")
	}
	return s.postRepo.Delete(ctx, id)
}
