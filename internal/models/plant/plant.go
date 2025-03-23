package plant

import (
	"context"
	"fmt"

	"time"

	"github.com/google/uuid"
)

type Plant struct {
	id            uuid.UUID
	name          string
	latinName     string
	description   string
	mainPhotoID   uuid.UUID
	photos        []PlantPhoto
	category      string
	specification PlantSpecification
	createdAt     time.Time
}

func (p *Plant) ID() uuid.UUID {
	return p.id
}

func (p *Plant) GetName() string {
	return p.name
}

func (p *Plant) GetDescription() string {
	return p.description
}

func (p *Plant) GetLatinName() string {
	return p.latinName
}

func (p *Plant) GetCategory() string {
	return p.category
}

func (p *Plant) GetSpecification() PlantSpecification {
	return p.specification
}

func (p *Plant) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Plant) MainPhotoID() uuid.UUID {
	return p.mainPhotoID
}

func (p *Plant) GetPhotos() []PlantPhoto {
	return p.photos
}

func (p *Plant) UpdateSpec(ps PlantSpecification) error {
	if err := ps.Validate(); err != nil {
		return err
	}
	p.specification = ps
	return nil
}

func (p *Plant) AddPhoto(photo *PlantPhoto) error {
	for _, v := range p.photos {
		if v.Compare(photo) {
			return fmt.Errorf("photo already exists")
		}
	}
	p.photos = append(p.photos, *photo)
	return nil
}

type PlantSpecification interface {
	Validate() error
}

func CreatePlant(id uuid.UUID,
	name, latinName, description string,
	mainPhotoID uuid.UUID,
	photos []PlantPhoto, category string,
	specification PlantSpecification,
	createdAt time.Time) (*Plant, error) {

	plant := &Plant{
		id:            id,
		name:          name,
		latinName:     latinName,
		description:   description,
		mainPhotoID:   mainPhotoID,
		photos:        photos,
		category:      category,
		specification: specification,
		createdAt:     createdAt,
	}
	if err := plant.Validate(); err != nil {
		return nil, err
	}
	return plant, nil
}

func NewPlant(name, latinName, description string,
	mainPhotoID uuid.UUID,
	photos []PlantPhoto, category string,
	specification PlantSpecification) (*Plant, error) {
	return CreatePlant(uuid.New(), name, latinName, description, mainPhotoID, photos, category, specification, time.Now())
}

func (p *Plant) Validate() error {
	if p.id == uuid.Nil {
		return fmt.Errorf("plant ID cannot be empty")
	}
	if p.name == "" {
		return fmt.Errorf("plant name cannot be empty")
	}
	if p.latinName == "" {
		return fmt.Errorf("plant latin name cannot be empty")
	}
	if p.description == "" {
		return fmt.Errorf("plant description cannot be empty")
	}
	if p.mainPhotoID == uuid.Nil {
		return fmt.Errorf("plant main photo ID cannot be empty")
	}
	if p.category == "" {
		return fmt.Errorf("plant category cannot be empty")
	}
	if p.specification == nil || p.specification.Validate() != nil {
		return fmt.Errorf("%v is not a valid specification", p.specification)
	}
	if p.createdAt.After(time.Now()) {
		return fmt.Errorf("plant creation date cannot be in the future: %v", p.createdAt)
	}
	if p.photos == nil {
		return fmt.Errorf("plant photos cannot be nil")
	}
	for _, p := range p.photos {
		if err := p.Validate(); err != nil {
			return fmt.Errorf("plant photo validation failed: %v", err)
		}
	}
	return nil
}

type PlantRepository interface {
	Create(ctx context.Context, plant *Plant) (*Plant, error)
	Update(ctx context.Context, plantID uuid.UUID, updateFn func(*Plant) (*Plant, error)) (*Plant, error)
	Delete(ctx context.Context, plantID uuid.UUID) error
	Get(ctx context.Context, plantID uuid.UUID) (*Plant, error)
}
