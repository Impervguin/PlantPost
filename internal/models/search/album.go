package search

import (
	"PlantSite/internal/models/album"

	"github.com/google/uuid"
)

type AlbumFilter interface {
	Filter(a *album.Album) bool
}

type AlbumOwnerFilter struct {
	OwnerID uuid.UUID
}

func NewAlbumOwnerFilter(ownerID uuid.UUID) *AlbumOwnerFilter {
	return &AlbumOwnerFilter{OwnerID: ownerID}
}

func (a *AlbumOwnerFilter) Filter(album *album.Album) bool {
	return album.GetOwnerID() == a.OwnerID
}
