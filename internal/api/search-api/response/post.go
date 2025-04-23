package response

type SearchPostItem struct {
	ID        string `json:"id" form:"id" binding:"required"`
	Title     string `json:"title" form:"title" binding:"required"`
	Content   string `json:"content" form:"content" binding:"required"`
	Tags      []string
	Photos    []SearchPostPhoto
	AuthorID  string `json:"author_id" form:"author_id" binding:"required"`
	UpdatedAt string `json:"updated_at" form:"updated_at" binding:"required"`
	CreatedAt string `json:"created_at" form:"created_at" binding:"required"`
}

type SearchPostPhoto struct {
	ID          string `json:"id" form:"id" binding:"required"`
	PlaceNumber int    `json:"place_number" form:"place_number" binding:"required"`
	Key         string `json:"key" form:"key" binding:"required"`
}

type SearchPostsResponse []SearchPostItem
