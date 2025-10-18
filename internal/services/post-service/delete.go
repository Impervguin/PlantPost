package postservice

import (
	"PlantSite/internal/models/auth"
	"PlantSite/internal/utils/logs"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *PostService) Delete(ctx context.Context, id uuid.UUID) error {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return auth.ErrNotAuthorized
	}
	logs.Debugf("PostService.Delete: got user=%s", user.ID())
	if !user.HasAuthorRights() {
		return auth.ErrNoAuthorRights
	}
	logs.Debugf("PostService.Delete: user %s has author rights", user.ID())
	if id == uuid.Nil {
		return fmt.Errorf("nil post")
	}
	logs.Debugf("PostService.Delete: got post ID %s", id)
	return s.postRepo.Delete(ctx, id)
}
