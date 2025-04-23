package request

import "github.com/google/uuid"

type CreateAlbumRequest struct {
	Name        string     `json:"name" form:"name" binding:"required"`
	Description string     `json:"description" form:"description" binding:"required"`
	PlantIDs    uuid.UUIDs `json:"plant_ids" form:"plant_ids" binding:"required"`
}

type GetAlbumRequest struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

type UpdateAlbumNameRequest struct {
	ID   uuid.UUID `uri:"id" binding:"required"`
	Name string    `json:"name" form:"name" binding:"required"`
}

type UpdateAlbumDescriptionRequest struct {
	ID          uuid.UUID `uri:"id" binding:"required"`
	Description string    `json:"description" form:"description" binding:"required"`
}

type AddPlantToAlbumRequest struct {
	ID      uuid.UUID `uri:"id" binding:"required"`
	PlantID uuid.UUID `json:"plant_id" form:"plant_id" binding:"required"`
}

type RemovePlantFromAlbumRequest struct {
	ID      uuid.UUID `uri:"id" binding:"required"`
	PlantID uuid.UUID `json:"plant_id" form:"plant_id" binding:"required"`
}

type DeleteAlbumRequest struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}
