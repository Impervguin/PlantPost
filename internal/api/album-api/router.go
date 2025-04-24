package albumapi

import (
	"PlantSite/internal/api/album-api/mapper"
	_ "PlantSite/internal/api/album-api/request"
	_ "PlantSite/internal/api/album-api/response"
	"PlantSite/internal/models/album"
	"PlantSite/internal/models/auth"
	albumservice "PlantSite/internal/services/album-service"
	"errors"
	"net/http"

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

// Album Create Handler
// @Summary Create an album
// @Description Creates a new album with the provided name and description
// @Tags album
// @Accept json
// @Param request body mapper.CreateAlbumRequest true "Create album request body"
// @Success 200  "Album created successfully"
// @Failure 400  "Bad Request - Invalid input or missing required fields"
// @Failure 401  "Unauthorized - Not authorized to create album"
// @Failure 500 "Internal Server Error - Failed to create album"
// @Router /album/create [post]
func (r *AlbumRouter) Create(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapCreateAlbumRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	// TODO: rewrite with CQRS in mind (or smth like that)
	alb, err := album.NewAlbum(req.Name, req.Description, req.PlantIDs, uuid.New())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	_, err = r.album.CreateAlbum(ctx, alb)
	if errors.Is(err, auth.ErrNotAuthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Error(err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// Get Album Handler
// @Summary Get album
// @Description Gets an album by ID
// @Tags album
// @Produce json
// @Param id path string true "Album ID"
// @Success 200  {object} response.GetAlbumResponse "Album fetch successfully"
// @Failure 400  "Bad Request - Invalid input"
// @Failure 401  "Unauthorized - Not authorized to get album"
// @Failure 500 "Internal Server Error - Failed to get album"
// @Router /album/get/{id} [get]
func (r *AlbumRouter) Get(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapGetAlbumRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	alb, err := r.album.GetAlbum(ctx, req.ID)
	if errors.Is(err, auth.ErrNotAuthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Error(err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	resp, err := mapper.MapGetAlbumResponse(alb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"album": resp})
}

// Update Album Name Handler
// @Summary Update album name
// @Description Updates the name of an album
// @Tags album
// @Accept json
// @Param id path string true "Album ID"
// @Param request body mapper.UpdateAlbumNameRequest true "Update album name request body"
// @Success 200  "Album name updated successfully"
// @Failure 400  "Bad Request - Invalid input or missing required fields"
// @Failure 401  "Unauthorized - Not authorized to update album name"
// @Failure 500 "Internal Server Error - Failed to update album name"
// @Router /album/name/{id} [put]
func (r *AlbumRouter) UpdateName(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapUpdateAlbumNameRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	err = r.album.UpdateAlbumName(ctx, req.ID, req.Name)
	if errors.Is(err, auth.ErrNotAuthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Error(err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// Update Album Description Handler
// @Summary Update album description
// @Description Updates the description of an album
// @Tags album
// @Accept json
// @Param id path string true "Album ID"
// @Param request body mapper.UpdateAlbumDescriptionRequest true "Update album description request body"
// @Success 200  "Album description updated successfully"
// @Failure 400  "Bad Request - Invalid input or missing required fields"
// @Failure 401  "Unauthorized - Not authorized to update album description"
// @Failure 500 "Internal Server Error - Failed to update album description"
// @Router /album/description/{id} [put]
func (r *AlbumRouter) UpdateDescription(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapUpdateAlbumDescriptionRequest(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	err = r.album.UpdateAlbumDescription(ctx, req.ID, req.Description)
	if errors.Is(err, auth.ErrNotAuthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Error(err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// Add Plant to Album Handler
// @Summary Add plant to album
// @Description Adds a plant to an album
// @Tags album
// @Accept json
// @Param id path string true "Album ID"
// @Param request body mapper.AddPlantToAlbumRequest true "Add plant to album request body"
// @Success 200  "Plant added to album successfully"
// @Failure 400  "Bad Request - Invalid input or missing required fields"
// @Failure 401  "Unauthorized - Not authorized to add plant to album"
// @Failure 500 "Internal Server Error - Failed to add plant to album"
// @Router /album/add/{id} [post]
func (r *AlbumRouter) AddPlantToAlbum(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapAddPlantToAlbumRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	err = r.album.AddPlantToAlbum(ctx, req.ID, req.PlantID)
	if errors.Is(err, auth.ErrNotAuthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Error(err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// Remove Plant from Album Handler
// @Summary Remove plant from album
// @Description Removes a plant from an album
// @Tags album
// @Accept json
// @Param id path string true "Album ID"
// @Param request body mapper.RemovePlantFromAlbumRequest true "Remove plant from album request body"
// @Success 200  "Plant removed from album successfully"
// @Failure 400  "Bad Request - Invalid input or missing required fields"
// @Failure 401  "Unauthorized - Not authorized to remove plant from album"
// @Failure 500 "Internal Server Error - Failed to remove plant from album"
// @Router /album/remove/{id} [delete]
func (r *AlbumRouter) RemovePlantFromAlbum(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapRemovePlantFromAlbumRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	err = r.album.RemovePlantFromAlbum(ctx, req.ID, req.PlantID)
	if errors.Is(err, auth.ErrNotAuthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Error(err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// Delete Album Handler
// @Summary Delete album
// @Description Deletes an album by ID
// @Tags album
// @Param id path string true "Album ID"
// @Success 200  "Album deleted successfully"
// @Failure 400  "Bad Request - Invalid input"
// @Failure 401  "Unauthorized - Not authorized to delete album"
// @Failure 500 "Internal Server Error - Failed to delete album"
// @Router /album/delete/{id} [delete]
func (r *AlbumRouter) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := mapper.MapDeleteAlbumRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	err = r.album.DeleteAlbum(ctx, req.ID)
	if errors.Is(err, auth.ErrNotAuthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Error(err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// List Albums Handler
// @Summary List albums
// @Description Lists all albums of authenticated user
// @Tags album
// @Produce json
// @Success 200  {object} response.ListAlbumsResponse "Albums fetch successfully"
// @Failure 401  "Unauthorized - Not authorized to list albums"
// @Failure 500 "Internal Server Error - Failed to list albums"
// @Router /album/list [get]
func (r *AlbumRouter) List(c *gin.Context) {
	ctx := c.Request.Context()

	albs, err := r.album.ListAlbums(ctx)
	if errors.Is(err, auth.ErrNotAuthorized) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Error(err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	resp, err := mapper.MapListAlbumsResponse(albs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"albums": resp})
}
