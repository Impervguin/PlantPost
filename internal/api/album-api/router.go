package albumapi

import (
	"PlantSite/internal/api/album-api/mapper"
	"PlantSite/internal/models/album"
	albumservice "PlantSite/internal/services/album-service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AlbumRouter struct {
	album *albumservice.AlbumService
}

func (r *AlbumRouter) Init(router *gin.RouterGroup, album *albumservice.AlbumService) {
	r.album = album
	gr := router.Group("/album")
	gr.POST("/create", r.Create)
	gr.GET("/get/:id", r.Get)
	gr.PUT("/name/:id", r.UpdateName)
	gr.PUT("/description/:id", r.UpdateDescription)
	gr.POST("/add/:id", r.AddPlantToAlbum)
	gr.DELETE("/remove/:id", r.RemovePlantFromAlbum)
	gr.DELETE("/delete/:id", r.Delete)
	gr.GET("/list", r.List)
}

func (r *AlbumRouter) Create(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapCreateAlbumRequest(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	// TODO: rewrite with CQRS in mind (or smth like that)
	alb, err := album.NewAlbum(req.Name, req.Description, req.PlantIDs, uuid.New())
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	_, err = r.album.CreateAlbum(ctx, alb)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{})
}

func (r *AlbumRouter) Get(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapGetAlbumRequest(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	alb, err := r.album.GetAlbum(ctx, req.ID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	resp, err := mapper.MapGetAlbumResponse(alb)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"album": resp})
}

func (r *AlbumRouter) UpdateName(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapUpdateAlbumNameRequest(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	err = r.album.UpdateAlbumName(ctx, req.ID, req.Name)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{})
}

func (r *AlbumRouter) UpdateDescription(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapUpdateAlbumDescriptionRequest(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	err = r.album.UpdateAlbumDescription(ctx, req.ID, req.Description)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{})
}

func (r *AlbumRouter) AddPlantToAlbum(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapAddPlantToAlbumRequest(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	err = r.album.AddPlantToAlbum(ctx, req.ID, req.PlantID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{})
}

func (r *AlbumRouter) RemovePlantFromAlbum(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapRemovePlantFromAlbumRequest(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	err = r.album.RemovePlantFromAlbum(ctx, req.ID, req.PlantID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{})
}

func (r *AlbumRouter) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapDeleteAlbumRequest(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	err = r.album.DeleteAlbum(ctx, req.ID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{})
}

func (r *AlbumRouter) List(c *gin.Context) {
	ctx := c.Request.Context()

	albs, err := r.album.ListAlbums(ctx)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	resp, err := mapper.MapListAlbumsResponse(albs)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{"albums": resp})
}
