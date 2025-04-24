package searchapi

import (
	"PlantSite/internal/api/search-api/mapper"
	_ "PlantSite/internal/api/search-api/request"
	_ "PlantSite/internal/api/search-api/response"
	"PlantSite/internal/models/search"
	searchservice "PlantSite/internal/services/search-service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchRouter struct {
	search *searchservice.SearchService
}

func (r *SearchRouter) Init(router *gin.RouterGroup, search *searchservice.SearchService) {
	r.search = search
	gr := router.Group("/search")
	gr.POST("/posts", r.SearchPosts)
	gr.POST("/plants", r.SearchPlants)
	gr.GET("/plant/:id", r.GetPlant)
	gr.GET("/post/:id", r.GetPost)
}

// @Summary Search posts with multiple filters
// @Description Search posts using an array of different filter types
// @Tags search
// @Accept json
// @Produce json
// @Param request body []mapper.SearchPostsItem true "Array of search filters"
// @Success 200 {object} response.SearchPostsResponse
// @Failure 400 "Invalid request format or missing required fields"
// @Failure 500 "Internal server error"
// @Router /search/posts [post]
func (r *SearchRouter) SearchPosts(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapSearchPostsRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	srchReq := search.NewPostSearch()
	for _, f := range req {
		filt, err := f.ToDomain()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Error(err)
			return
		}
		srchReq.AddFilter(filt)
	}

	posts, err := r.search.SearchPosts(ctx, srchReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	resp := mapper.MapSearchPostsResponse(posts)
	if resp == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": resp})
}

// @Summary Search plants with multiple filters
// @Description Search plants using an array of different filter types
// @Tags search
// @Accept json
// @Produce json
// @Param request body []mapper.SearchPlantsItem true "Array of search filters"
// @Success 200 {object} response.SearchPlantResponse
// @Failure 400 "Invalid request format or missing required fields"
// @Failure 500 "Internal server error"
// @Router /search/plants [post]
func (r *SearchRouter) SearchPlants(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapSearchPlantsRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	srchReq := search.NewPlantSearch()
	for _, f := range req {
		filt, err := f.ToDomain()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Error(err)
			return
		}
		srchReq.AddFilter(filt)
	}

	plants, err := r.search.SearchPlants(ctx, srchReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	resp, err := mapper.MapSearchPlantsResponse(plants)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"plants": resp})
}

// @Summary Get post
// @Description Gets a post by ID
// @Tags search
// @Produce json
// @Param id path string true "Post ID"
// @Success 200  {object} response.GetPostResponse "Post fetch successfully"
// @Failure 400  "Bad Request - Invalid input"
// @Failure 401  "Unauthorized - Not authorized to get post"
// @Failure 403  "Forbidden - Does not have author rights to get post"
// @Failure 500 "Internal Server Error - Failed to get post"
// @Router /search/post/{id} [get]
func (r *SearchRouter) GetPost(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapGetPostRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	pst, err := r.search.GetPost(ctx, req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	resp := mapper.MapGetPostResponse(pst)
	if resp == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": resp})
}

// @Summary Get plant
// @Description Gets a plant by ID
// @Tags search
// @Produce json
// @Param id path string true "Plant ID"
// @Success 200  {object} response.GetPlantResponse "Plant fetch successfully"
// @Failure 400  "Bad Request - Invalid input"
// @Failure 401  "Unauthorized - Not authorized to get plant"
// @Failure 403  "Forbidden - Does not have author rights to get plant"
// @Failure 500 "Internal Server Error - Failed to get plant"
// @Router /search/plant/{id} [get]
func (r *SearchRouter) GetPlant(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapGetPlantRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	pl, err := r.search.GetPlantByID(ctx, req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	resp, err := mapper.MapGetPlantResponse(pl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"plant": resp})
}
