package mapper

import (
	"PlantSite/internal/api/post-api/response"
	postservice "PlantSite/internal/services/post-service"
)

const timeFormat = "2006-01-02 15:04:05"

func MapGetPostResponse(pst *postservice.GetPost) *response.GetPostResponse {
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

func MapGetPostPhotos(photos []postservice.GetPostPhoto) []response.GetPostPhoto {
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
