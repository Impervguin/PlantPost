package searchapi

import (
	"PlantSite/internal/api/search-api/mapper"
	"PlantSite/internal/models/search"
	searchservice "PlantSite/internal/services/search-service"
	"fmt"

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
}

func (r *SearchRouter) SearchPosts(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapSearchPostsRequest(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	srchReq := search.NewPostSearch()
	for _, f := range req {
		filt, err := f.ToDomain()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			c.Error(err)
			return
		}
		srchReq.AddFilter(filt)
	}

	posts, err := r.search.SearchPosts(ctx, srchReq)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	resp := mapper.MapSearchPostsResponse(posts)
	if resp == nil {
		c.JSON(404, gin.H{"error": "post not found"})
		return
	}

	c.JSON(200, gin.H{"posts": resp})
}

func (r *SearchRouter) SearchPlants(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapSearchPlantsRequest(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	srchReq := search.NewPlantSearch()
	for _, f := range req {
		filt, err := f.ToDomain()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			c.Error(err)
			return
		}
		srchReq.AddFilter(filt)
	}
	fmt.Println("test here")
	srchReq.Iterate(func(f search.PlantFilter) error {
		fmt.Println(f.Identifier())
		return nil
	})

	plants, err := r.search.SearchPlants(ctx, srchReq)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	resp, err := mapper.MapSearchPlantsResponse(plants)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"plants": resp})
}
