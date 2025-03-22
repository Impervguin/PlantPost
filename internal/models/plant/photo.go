package plant

import (
	"fmt"

	"github.com/google/uuid"
)

type PlantPhoto struct {
	id          uuid.UUID
	fileID      uuid.UUID
	description string
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
