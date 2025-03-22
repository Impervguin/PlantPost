package album

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Album struct {
	id          uuid.UUID
	name        string
	description string
	plantIDs    uuid.UUIDs
	ownerID     uuid.UUID
	createdAt   time.Time
	updatedAt   time.Time
}

func CreateAlbum(id uuid.UUID,
	name, description string,
	plantIDs uuid.UUIDs,
	ownerID uuid.UUID,
	createdAt time.Time,
	updatedAt time.Time) (*Album, error) {
	album := &Album{
		id:          id,
		name:        name,
		description: description,
		plantIDs:    plantIDs,
		ownerID:     ownerID,
		createdAt:   createdAt,
	}
	if err := album.Validate(); err != nil {
		return nil, err
	}
	return album, nil
}

func NewAlbum(name, description string,
	plantIDs uuid.UUIDs,
	ownerID uuid.UUID) (*Album, error) {
	return CreateAlbum(uuid.New(), name, description, plantIDs, ownerID, time.Now(), time.Now())
}

func (album *Album) Validate() error {
	if album.id == uuid.Nil {
		return fmt.Errorf("album id cannot be nil")
	}
	if album.name == "" {
		return fmt.Errorf("album name cannot be empty")
	}

	if album.ownerID == uuid.Nil {
		return fmt.Errorf("album owner id cannot be nil")
	}
	if album.updatedAt.After(time.Now()) {
		return fmt.Errorf("can't be updated in future %v", album.updatedAt)
	}
	if album.createdAt.After(album.updatedAt) {
		return fmt.Errorf("can't be created after update %v", album.createdAt)
	}
	if album.plantIDs == nil {
		return fmt.Errorf("album plant ids cannot be nil")
	}
	for _, p := range album.plantIDs {
		if p == uuid.Nil {
			return fmt.Errorf("plant id cannot be nil")
		}
	}

	return nil
}

func (album Album) GetOwnerID() uuid.UUID {
	return album.ownerID
}

type AlbumRepository interface {
	Create(alb *Album) (*Album, error)
	Update(id uuid.UUID, updateFn func(*Album) (*Album, error)) (*Album, error)
	Delete(id uuid.UUID) error
	Get(id uuid.UUID) (*Album, error)
}
