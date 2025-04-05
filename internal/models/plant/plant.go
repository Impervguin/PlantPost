package plant

import (
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
	photos        PlantPhotos
	category      string
	specification PlantSpecification
	createdAt     time.Time
	updatedAt     time.Time
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

func (p *Plant) UpdatedAt() time.Time {
	return p.updatedAt
}

func (p *Plant) MainPhotoID() uuid.UUID {
	return p.mainPhotoID
}

func (p *Plant) GetPhotos() PlantPhotos {
	return p.photos
}

func (p *Plant) UpdateSpec(ps PlantSpecification) error {
	if err := ps.Validate(); err != nil {
		return err
	}
	p.specification = ps
	p.updatedAt = time.Now()
	return nil
}

func (p *Plant) AddPhoto(photo *PlantPhoto) error {
	err := p.photos.Add(photo)
	if err != nil {
		return err
	}
	p.updatedAt = time.Now()
	return nil
}

type PlantSpecification interface {
	Validate() error
	Category() string
}

func CreatePlant(id uuid.UUID,
	name, latinName, description string,
	mainPhotoID uuid.UUID,
	photos PlantPhotos, category string,
	specification PlantSpecification,
	createdAt time.Time, updatedAt time.Time) (*Plant, error) {

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
		updatedAt:     updatedAt,
	}
	if err := plant.Validate(); err != nil {
		return nil, err
	}
	return plant, nil
}

func NewPlant(name, latinName, description string,
	mainPhotoID uuid.UUID,
	photos PlantPhotos, category string,
	specification PlantSpecification) (*Plant, error) {
	return CreatePlant(uuid.New(), name, latinName, description, mainPhotoID, photos, category, specification, time.Now(), time.Now())
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
	if p.updatedAt.After(time.Now()) {
		return fmt.Errorf("plant update date cannot be in the future: %v", p.createdAt)
	}
	if p.createdAt.After(p.updatedAt) {
		return fmt.Errorf("plant creation date cannot be after update date: %v %v", p.createdAt, p.updatedAt)
	}
	if err := p.photos.Validate(); err != nil {
		return err
	}
	return nil
}
