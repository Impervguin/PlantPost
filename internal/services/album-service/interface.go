package albumservice

import (
	"PlantSite/internal/models/album"
	"context"

	"github.com/google/uuid"
)

type AlbumServiceContract interface {
	CreateAlbum(ctx context.Context, alb *album.Album) (*album.Album, error)
	GetAlbum(ctx context.Context, id uuid.UUID) (*album.Album, error)
	UpdateAlbumName(ctx context.Context, id uuid.UUID, name string) error
	UpdateAlbumDescription(ctx context.Context, id uuid.UUID, description string) error
	AddPlantToAlbum(ctx context.Context, id uuid.UUID, plantID uuid.UUID) error
	RemovePlantFromAlbum(ctx context.Context, id uuid.UUID, plantID uuid.UUID) error
	DeleteAlbum(ctx context.Context, id uuid.UUID) error
	ListAlbums(ctx context.Context) ([]*album.Album, error)
}
