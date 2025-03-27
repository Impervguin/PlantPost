package plant

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
)

type PlantPhoto struct {
	id          uuid.UUID
	fileID      uuid.UUID
	description string
}

type PlantPhotos struct {
	photos []PlantPhoto
}

func CreatePlantPhoto(id, fileID uuid.UUID, description string) (*PlantPhoto, error) {
	pphoto := &PlantPhoto{
		id:          id,
		fileID:      fileID,
		description: description,
	}
	if err := pphoto.Validate(); err != nil {
		return nil, err
	}
	return pphoto, nil
}

func (pphoto *PlantPhoto) Compare(other *PlantPhoto) bool {
	if pphoto == nil || other == nil {
		return false
	}
	return pphoto.id == other.id || pphoto.fileID == other.fileID
}

func (pphoto *PlantPhoto) FileID() uuid.UUID {
	return pphoto.fileID
}

func (pphoto *PlantPhoto) Description() string {
	return pphoto.description
}

func (pphoto *PlantPhoto) ID() uuid.UUID {
	return pphoto.id
}

func NewPlantPhoto(fileID uuid.UUID, description string) (*PlantPhoto, error) {
	return CreatePlantPhoto(uuid.New(), fileID, description)
}

func (pphoto *PlantPhoto) Validate() error {
	if pphoto.id == uuid.Nil {
		return fmt.Errorf("plant photo ID cannot be empty")
	}
	if pphoto.fileID == uuid.Nil {
		return fmt.Errorf("plant photo file ID cannot be empty")
	}
	return nil
}

func NewPlantPhotos() *PlantPhotos {
	return &PlantPhotos{photos: make([]PlantPhoto, 0)}
}

func (pp *PlantPhotos) Add(photo *PlantPhoto) error {
	if slices.ContainsFunc(pp.photos, func(e PlantPhoto) bool {
		return e.Compare(photo)
	}) {
		return fmt.Errorf("photo already exists")
	}
	pp.photos = append(pp.photos, *photo)
	return nil
}

func (pp *PlantPhotos) Remove(photo *PlantPhoto) error {
	index := slices.IndexFunc(pp.photos, func(e PlantPhoto) bool {
		return e.Compare(photo)
	})
	if index == -1 {
		return fmt.Errorf("photo not found")
	}
	pp.photos = append(pp.photos[:index], pp.photos[index+1:]...)
	return nil
}

func (pp PlantPhotos) Iterate(iterFunc func(e PlantPhoto) error) error {
	for _, photo := range pp.photos {
		if err := iterFunc(photo); err != nil {
			return err
		}
	}
	return nil
}

func (pp *PlantPhotos) IterateUpdate(iterFunc func(e *PlantPhoto) error) error {
	for i := range pp.photos {
		photoCopy := pp.photos[i]
		if err := iterFunc(&photoCopy); err != nil {
			return err
		}
		if err := photoCopy.Validate(); err != nil {
			return err
		}
		pp.photos[i] = photoCopy
	}
	return nil
}

func (pp PlantPhotos) Len() int { return len(pp.photos) }

func (pp *PlantPhotos) Validate() error {
	return nil
}
