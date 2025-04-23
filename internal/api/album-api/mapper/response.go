package mapper

import (
	"PlantSite/internal/api/album-api/response"
	"PlantSite/internal/models/album"
)

const timeFormat = "2006-01-02 15:04:05"

func MapGetAlbumResponse(alb *album.Album) (*response.GetAlbumResponse, error) {
	if alb == nil {
		return nil, nil
	}
	plantIDs := make([]string, 0, len(alb.PlantIDs()))
	for _, id := range alb.PlantIDs() {
		plantIDs = append(plantIDs, id.String())
	}
	return &response.GetAlbumResponse{
		ID:          alb.ID().String(),
		Name:        alb.Name(),
		Description: alb.Description(),
		PlantIDs:    plantIDs,
		CreatedAt:   alb.CreatedAt().Format(timeFormat),
		UpdatedAt:   alb.UpdatedAt().Format(timeFormat),
	}, nil
}

func MapListAlbumsResponse(albs []*album.Album) (*response.ListAlbumsResponse, error) {
	if albs == nil {
		return &response.ListAlbumsResponse{}, nil
	}
	resp := make(response.ListAlbumsResponse, 0, len(albs))
	for _, alb := range albs {
		plantIDs := make([]string, 0, len(alb.PlantIDs()))
		for _, id := range alb.PlantIDs() {
			plantIDs = append(plantIDs, id.String())
		}
		resp = append(resp, response.ListAlbum{
			ID:          alb.ID().String(),
			Name:        alb.Name(),
			Description: alb.Description(),
			PlantIDs:    plantIDs,
			CreatedAt:   alb.CreatedAt().Format(timeFormat),
			UpdatedAt:   alb.UpdatedAt().Format(timeFormat),
		})
	}
	return &resp, nil
}
