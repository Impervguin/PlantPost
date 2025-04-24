package mapper

import (
	"PlantSite/internal/api/plant-api/spec"
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
		respSpec, err := spec.MapSpecification(p.Specification)
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

func MapGetPlantResponse(pl *searchservice.GetPlant) (*response.GetPlantResponse, error) {
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

func MapGetPostResponse(pst *searchservice.GetPost) *response.GetPostResponse {
	if pst == nil {
		return nil
	}
	return &response.GetPostResponse{
		ID:          pst.ID,
		Title:       pst.Title,
		Content:     pst.Content.Text,
		ContentType: string(pst.Content.ContentType),
		Tags:        pst.Tags,
		Photos:      MapGetPostPhotos(pst.Photos),

		AuthorID:  pst.AuthorID,
		UpdatedAt: pst.UpdatedAt.Format(timeFormat),
		CreatedAt: pst.CreatedAt.Format(timeFormat),
	}
}

func MapGetPostPhotos(photos []searchservice.GetPostPhoto) []response.GetPostPhoto {
	if photos == nil {
		return nil
	}
	res := make([]response.GetPostPhoto, 0)
	for _, photo := range photos {
		res = append(res, response.GetPostPhoto{
			ID:          photo.ID,
			PlaceNumber: photo.PlaceNumber,
			Key:         photo.File.URL,
		})
	}
	return res
}
