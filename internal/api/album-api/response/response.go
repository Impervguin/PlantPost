package response

type GetAlbumResponse struct {
	ID          string `json:"id" form:"id" binding:"required"`
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description" form:"description" binding:"required"`
	PlantIDs    []string
	CreatedAt   string `json:"created_at" form:"created_at" binding:"required"`
	UpdatedAt   string `json:"updated_at" form:"updated_at" binding:"required"`
}

type ListAlbum struct {
	ID          string `json:"id" form:"id" binding:"required"`
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description" form:"description" binding:"required"`
	PlantIDs    []string
	CreatedAt   string `json:"created_at" form:"created_at" binding:"required"`
	UpdatedAt   string `json:"updated_at" form:"updated_at" binding:"required"`
}

type ListAlbumsResponse []ListAlbum
