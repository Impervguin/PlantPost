package albumservice

import (
	"PlantSite/internal/models/album"
	authservice "PlantSite/internal/services/auth-service"
	"context"

	"github.com/google/uuid"
)

type AlbumService struct {
	albumRepository album.AlbumRepository
	auth            *authservice.AuthService
}

func NewAlbumService(repo album.AlbumRepository, auth *authservice.AuthService) *AlbumService {
	return &AlbumService{
		albumRepository: repo,
		auth:            auth,
	}
}

func (s *AlbumService) CreateAlbum(ctx context.Context, alb *album.Album) (*album.Album, error) {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return nil, ErrNotAuthorized
	}
	if !user.HasMemberRights() {
		return nil, ErrNotMember
	}
	if alb.GetOwnerID() != user.ID() {
		return nil, ErrNotOwner
	}
	alb, err := s.albumRepository.Create(ctx, alb)
	if err != nil {
		return nil, Wrap(err)
	}
	return alb, nil
}

func (s *AlbumService) GetAlbum(ctx context.Context, id uuid.UUID) (*album.Album, error) {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return nil, ErrNotAuthorized
	}
	if !user.HasMemberRights() {
		return nil, ErrNotMember
	}
	alb, err := s.albumRepository.Get(ctx, id)
	if err != nil {
		return nil, Wrap(err)
	}
	if alb.GetOwnerID() != user.ID() {
		return nil, ErrNotOwner
	}
	return alb, nil
}

func (s *AlbumService) UpdateAlbumName(ctx context.Context, id uuid.UUID, name string) error {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return ErrNotAuthorized
	}
	if !user.HasMemberRights() {
		return ErrNotMember
	}
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
	return nil
}

func (s *AlbumService) UpdateAlbumDescription(ctx context.Context, id uuid.UUID, description string) error {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return ErrNotAuthorized
	}
	if !user.HasMemberRights() {
		return ErrNotMember
	}
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
	return nil
}

func (s *AlbumService) AddPlantToAlbum(ctx context.Context, id uuid.UUID, plantID uuid.UUID) error {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return ErrNotAuthorized
	}
	if !user.HasMemberRights() {
		return ErrNotMember
	}
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
	return nil
}

func (s *AlbumService) RemovePlantFromAlbum(ctx context.Context, id uuid.UUID, plantID uuid.UUID) error {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return ErrNotAuthorized
	}
	if !user.HasMemberRights() {
		return ErrNotMember
	}
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
	return nil
}

func (s *AlbumService) DeleteAlbum(ctx context.Context, id uuid.UUID) error {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return ErrNotAuthorized
	}
	if !user.HasMemberRights() {
		return ErrNotMember
	}
	alb, err := s.albumRepository.Get(ctx, id)
	if err != nil {
		return Wrap(err)
	}
	if alb.GetOwnerID() != user.ID() {
		return ErrNotOwner
	}
	err = s.albumRepository.Delete(ctx, id)
	if err != nil {
		return Wrap(err)
	}
	return nil
}

func (s *AlbumService) ListAlbums(ctx context.Context) ([]*album.Album, error) {
	user := s.auth.UserFromContext(ctx)
	if user == nil {
		return nil, ErrNotAuthorized
	}
	if !user.HasMemberRights() {
		return nil, ErrNotMember
	}
	albs, err := s.albumRepository.List(ctx, user.ID())
	if err != nil {
		return nil, Wrap(err)
	}
	return albs, nil
}
