package mapper

import (
	"PlantSite/internal/api/plant-api/response"
	"PlantSite/internal/api/plant-api/spec"
	plantservice "PlantSite/internal/services/plant-service"
	"fmt"
)

var timeFormat = "2006-01-02 15:04:05"

func MapGetPlantResponse(pl *plantservice.GetPlant) (*response.GetPlantResponse, error) {
	if pl == nil {
		return nil, nil
	}
	spec, err := spec.MapSpecification(pl.Specification)
	if err != nil {
		return nil, fmt.Errorf("can't map specification to mapper: %w", err)
	}

	photos := make([]response.GetPlantPhoto, 0, len(pl.Photos))
	for _, photo := range pl.Photos {
		photos = append(photos, response.GetPlantPhoto{
			ID:          photo.ID.String(),
			Key:         photo.File.URL,
			Description: photo.Description,
		})
	}

	return &response.GetPlantResponse{
		ID:            pl.ID.String(),
		Name:          pl.Name,
		LatinName:     pl.LatinName,
		Description:   pl.Description,
		MainPhotoKey:  pl.MainPhoto.URL,
		Photos:        photos,
		Category:      pl.Category,
		Specification: spec,
		CreatedAt:     pl.CreatedAt.Format(timeFormat),
	}, nil
}
