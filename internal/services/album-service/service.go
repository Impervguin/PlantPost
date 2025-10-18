package albumservice

import (
	"PlantSite/internal/models/album"
	"PlantSite/internal/models/auth"
	authservice "PlantSite/internal/services/auth-service"
	"PlantSite/internal/utils/logs"
	"context"

	"github.com/google/uuid"
)

type AlbumService struct {
	albumRepository album.AlbumRepository
	auth            authservice.AuthServiceContract
}

func NewAlbumService(repo album.AlbumRepository, auth authservice.AuthServiceContract) *AlbumService {
	return &AlbumService{
		albumRepository: repo,
		auth:            auth,
	}
}

func (s *AlbumService) CreateAlbum(ctx context.Context, alb *album.Album) (*album.Album, error) {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return nil, auth.ErrNotAuthorized
	}
	logs.Debugf("AlbumService.CreateAlbum: got user=%s", user.ID())
	if !user.HasMemberRights() {
		return nil, auth.ErrNoMemberRights
	}
	logs.Debugf("AlbumService.CreateAlbum: user %s has member rights", user.ID())
	ownerAlb, err := album.CreateAlbum(
		alb.ID(),
		alb.Name(),
		alb.Description(),
		alb.PlantIDs(),
		user.ID(),
		alb.CreatedAt(),
		alb.UpdatedAt(),
	)
	if err != nil {
		return nil, Wrap(err)
	}
	logs.Debugf("AlbumService.CreateAlbum: created album %s for user %s", alb.ID(), user.ID())
	alb, err = s.albumRepository.Create(ctx, ownerAlb)
	if err != nil {
		return nil, Wrap(err)
	}
	logs.Debugf("AlbumService.CreateAlbum: saved album %s for user %s", alb.ID(), user.ID())
	return alb, nil
}

func (s *AlbumService) GetAlbum(ctx context.Context, id uuid.UUID) (*album.Album, error) {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return nil, auth.ErrNotAuthorized
	}
	logs.Debugf("AlbumService.GetAlbum: got user=%s", user.ID())
	if !user.HasMemberRights() {
		return nil, auth.ErrNoMemberRights
	}
	logs.Debugf("AlbumService.GetAlbum: user %s has member rights", user.ID())
	alb, err := s.albumRepository.Get(ctx, id)
	if err != nil {
		return nil, Wrap(err)
	}
	logs.Debugf("AlbumService.GetAlbum: got album %s for user %s", alb.ID(), user.ID())
	if alb.GetOwnerID() != user.ID() {
		return nil, ErrNotOwner
	}
	logs.Debugf("AlbumService.GetAlbum: album %s is owned by user %s", alb.ID(), user.ID())
	return alb, nil
}

func (s *AlbumService) UpdateAlbumName(ctx context.Context, id uuid.UUID, name string) error {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return auth.ErrNotAuthorized
	}
	logs.Debugf("AlbumService.UpdateAlbumName: got user=%s", user.ID())
	if !user.HasMemberRights() {
		return auth.ErrNoMemberRights
	}
	logs.Debugf("AlbumService.UpdateAlbumName: user %s has member rights", user.ID())
	_, err := s.albumRepository.Update(ctx, id, func(a *album.Album) (*album.Album, error) {
		if a.GetOwnerID() != user.ID() {
			return nil, ErrNotOwner
		}
		err := a.UpdateName(name)
		return a, err
	})
	if err != nil {
		return Wrap(err)
	}
	logs.Debugf("AlbumService.UpdateAlbumName: updated album %s name to %s for user %s", id, name, user.ID())
	return nil
}

func (s *AlbumService) UpdateAlbumDescription(ctx context.Context, id uuid.UUID, description string) error {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return auth.ErrNotAuthorized
	}
	logs.Debugf("AlbumService.UpdateAlbumDescription: updating album %s description for user %s", id, user.ID())
	if !user.HasMemberRights() {
		return auth.ErrNoMemberRights
	}
	logs.Debugf("AlbumService.UpdateAlbumDescription: updating album %s description for user %s", id, user.ID())
	_, err := s.albumRepository.Update(ctx, id, func(a *album.Album) (*album.Album, error) {
		if a.GetOwnerID() != user.ID() {
			return nil, ErrNotOwner
		}
		err := a.UpdateDescription(description)
		return a, err
	})
	if err != nil {
		return Wrap(err)
	}
	logs.Debugf("AlbumService.UpdateAlbumDescription: saved updated album %s description to %s for user %s", id, description, user.ID())
	return nil
}

func (s *AlbumService) AddPlantToAlbum(ctx context.Context, id uuid.UUID, plantID uuid.UUID) error {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return auth.ErrNotAuthorized
	}
	logs.Debugf("AlbumService.AddPlantToAlbum: got user=%s", user.ID())
	if !user.HasMemberRights() {
		return auth.ErrNoMemberRights
	}
	logs.Debugf("AlbumService.AddPlantToAlbum: user %s has member rights", user.ID())
	_, err := s.albumRepository.Update(ctx, id, func(a *album.Album) (*album.Album, error) {
		if a.GetOwnerID() != user.ID() {
			return nil, ErrNotOwner
		}
		err := a.AddPlant(plantID)
		return a, err
	})
	if err != nil {
		return Wrap(err)
	}
	logs.Debugf("AlbumService.AddPlantToAlbum: saved updated album %s with plant %s for user %s", id, plantID, user.ID())
	return nil
}

func (s *AlbumService) RemovePlantFromAlbum(ctx context.Context, id uuid.UUID, plantID uuid.UUID) error {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return auth.ErrNotAuthorized
	}
	logs.Debugf("AlbumService.RemovePlantFromAlbum: got user=%s", user.ID())
	if !user.HasMemberRights() {
		return auth.ErrNoMemberRights
	}
	logs.Debugf("AlbumService.RemovePlantFromAlbum: user %s has member rights", user.ID())
	_, err := s.albumRepository.Update(ctx, id, func(a *album.Album) (*album.Album, error) {
		if a.GetOwnerID() != user.ID() {
			return nil, ErrNotOwner
		}
		err := a.RemovePlant(plantID)
		return a, err
	})
	if err != nil {
		return Wrap(err)
	}
	logs.Debugf("AlbumService.RemovePlantFromAlbum: saved updated album %s with plant %s for user %s", id, plantID, user.ID())
	return nil
}

func (s *AlbumService) DeleteAlbum(ctx context.Context, id uuid.UUID) error {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return auth.ErrNotAuthorized
	}
	logs.Debugf("AlbumService.DeleteAlbum: got user=%s", user.ID())
	if !user.HasMemberRights() {
		return auth.ErrNoMemberRights
	}
	logs.Debugf("AlbumService.DeleteAlbum: user %s has member rights", user.ID())
	alb, err := s.albumRepository.Get(ctx, id)
	if err != nil {
		return Wrap(err)
	}
	logs.Debugf("AlbumService.DeleteAlbum: got album %s for user %s", alb.ID(), user.ID())
	if alb.GetOwnerID() != user.ID() {
		return ErrNotOwner
	}
	logs.Debugf("AlbumService.DeleteAlbum: album %s is owned by user %s", alb.ID(), user.ID())
	err = s.albumRepository.Delete(ctx, id)
	if err != nil {
		return Wrap(err)
	}
	logs.Debugf("AlbumService.DeleteAlbum: deleted album %s for user %s", id, user.ID())
	return nil
}

func (s *AlbumService) ListAlbums(ctx context.Context) ([]*album.Album, error) {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return nil, auth.ErrNotAuthorized
	}
	logs.Debugf("AlbumService.ListAlbums: got user=%s", user.ID())
	if !user.HasMemberRights() {
		return nil, auth.ErrNoMemberRights
	}
	logs.Debugf("AlbumService.ListAlbums: user %s has member rights", user.ID())
	albs, err := s.albumRepository.List(ctx, user.ID())
	if err != nil {
		return nil, Wrap(err)
	}
	logs.Debugf("AlbumService.ListAlbums: got %d albums for user %s", len(albs), user.ID())
	return albs, nil
}
