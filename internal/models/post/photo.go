package post

import (
	"fmt"
	"slices"
	"sort"

	"github.com/google/uuid"
)

type PostPhoto struct {
	id          uuid.UUID
	placeNumber int
	fileID      uuid.UUID
}

func CreatePostPhoto(id uuid.UUID, fileID uuid.UUID, placeNumber int) (*PostPhoto, error) {
	postPhoto := &PostPhoto{
		id:          id,
		placeNumber: placeNumber,
		fileID:      fileID,
	}
	return postPhoto, nil
}

func NewPostPhoto(fileID uuid.UUID, placeNumber int) (*PostPhoto, error) {
	return CreatePostPhoto(uuid.New(), fileID, placeNumber)
}

func (p *PostPhoto) ID() uuid.UUID {
	return p.id
}

func (p *PostPhoto) PlaceNumber() int {
	return p.placeNumber
}

func (p *PostPhoto) FileID() uuid.UUID {
	return p.fileID
}

func (p *PostPhoto) Validate() error {
	if p.id == uuid.Nil {
		return fmt.Errorf("post photo ID cannot be empty")
	}
	if p.fileID == uuid.Nil {
		return fmt.Errorf("photo file ID cannot be empty")
	}
	return nil
}

type PostPhotos struct {
	photos []PostPhoto
}

func NewPostPhotos() *PostPhotos {
	return &PostPhotos{}
}

func (pp *PostPhotos) Validate() error {
	photoSet := make(map[int]struct{})
	for _, photo := range pp.photos {
		if photo.PlaceNumber() < 1 || photo.PlaceNumber() > MaximumPhotoPerPostCount {
			return fmt.Errorf("invalid photo place number: %d", photo.placeNumber)
		}
		if _, exists := photoSet[photo.PlaceNumber()]; exists {
			return fmt.Errorf("duplicate photo place number: %d", photo.placeNumber)
		}
		photoSet[photo.PlaceNumber()] = struct{}{}
		if photo.FileID() == uuid.Nil {
			return fmt.Errorf("photo file ID cannot be empty")
		}
	}
	return nil
}

func (pp *PostPhotos) RebalancePositions() {

	// Sort photos by place number
	sort.Slice(pp.photos, func(i, j int) bool {
		return pp.photos[i].PlaceNumber() < pp.photos[j].PlaceNumber()
	})

	// Update place numbers to consecutive integers starting from 1
	for i := range pp.photos {
		pp.photos[i].placeNumber = i + 1
	}
}

func (pp *PostPhotos) Add(photo *PostPhoto) error {
	if len(pp.photos) >= MaximumPhotoPerPostCount {
		return fmt.Errorf("maximum number of photos exceeded")
	}
	if err := photo.Validate(); err != nil {
		return err
	}
	if slices.ContainsFunc(pp.photos, func(ph PostPhoto) bool { return photo.ID() == ph.ID() || photo.FileID() == ph.FileID() }) {
		return fmt.Errorf("photo already exists")
	}
	for i := range pp.photos {
		if pp.photos[i].PlaceNumber() >= photo.PlaceNumber() {
			pp.photos[i].placeNumber++
		}
	}
	pp.photos = append(pp.photos, *photo)

	pp.RebalancePositions()
	return nil
}

func (pp *PostPhotos) Remove(photoID uuid.UUID) error {
	for i, ph := range pp.photos {
		if ph.ID() == photoID {
			for j := i + 1; j < len(pp.photos); j++ {
				pp.photos[j].placeNumber--
			}
			pp.photos = append(pp.photos[:i], pp.photos[i+1:]...)
			pp.RebalancePositions()
			return nil
		}
	}
	return fmt.Errorf("photo not found")
}

func (pp PostPhotos) List() []PostPhoto {
	ppcopy := make([]PostPhoto, len(pp.photos))
	copy(ppcopy, pp.photos)
	return ppcopy
}

func (pp PostPhotos) Len() int {
	return len(pp.photos)
}
