package plant

import "github.com/google/uuid"

type PlantCategory struct {
	ID          uuid.UUID
	Name        string
	MainPhotoID uuid.UUID
}
