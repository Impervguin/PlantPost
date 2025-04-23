package mapper

import (
	plantspec "PlantSite/internal/api/plant-spec"
	"PlantSite/internal/api/search-api/response"
	searchservice "PlantSite/internal/services/search-service"
	"fmt"
)

const timeFormat = "2006-01-02 15:04:05"

func MapSearchPostsResponse(posts []*searchservice.SearchPost) response.SearchPostsResponse {
	if posts == nil {
		return nil
	}
	resp := make(response.SearchPostsResponse, 0, len(posts))
	for _, p := range posts {
		photos := make([]response.SearchPostPhoto, 0)
		for _, photo := range p.Photos {
			photos = append(photos, response.SearchPostPhoto{
				ID:          photo.ID.String(),
				PlaceNumber: photo.PlaceNumber,
				Key:         photo.File.URL,
			})
		}
		resp = append(resp, response.SearchPostItem{
			ID:        p.ID.String(),
			Title:     p.Title,
			Content:   p.Content.Text,
			Tags:      p.Tags,
			Photos:    photos,
			AuthorID:  p.AuthorID.String(),
			UpdatedAt: p.UpdatedAt.Format(timeFormat),
			CreatedAt: p.CreatedAt.Format(timeFormat),
		})
	}
	return resp
}

func MapSearchPlantsResponse(plants []*searchservice.SearchPlant) (response.SearchPlantResponse, error) {
	if plants == nil {
		return nil, fmt.Errorf("plants not found")
	}
	resp := make(response.SearchPlantResponse, 0, len(plants))

	for _, p := range plants {
		respSpec, err := plantspec.MapSpecification(p.Specification)
		if err != nil {
			return nil, fmt.Errorf("can't map specification to mapper: %w", err)
		}
		resp = append(resp, response.SearchPlantItem{
			ID:            p.ID.String(),
			Name:          p.Name,
			LatinName:     p.LatinName,
			Description:   p.Description,
			MainPhotoKey:  p.MainPhoto.URL,
			Category:      p.Category,
			Specification: respSpec,
			CreatedAt:     p.CreatedAt.Format(timeFormat),
		})
	}
	return resp, nil
}
