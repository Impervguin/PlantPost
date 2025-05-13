package view

import (
	"PlantSite/internal/models/album"
	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/search"
	albumservice "PlantSite/internal/services/album-service"
	searchservice "PlantSite/internal/services/search-service"
	"PlantSite/internal/view/components"
	"PlantSite/internal/view/gintemplrenderer"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *ViewRouter) AlbumsCreateHandler(c *gin.Context) {
	ctx := c.Request.Context()
	user := r.auth.UserFromContext(ctx)

	if !user.HasMemberRights() {
		c.Redirect(http.StatusFound, "/view/albums")
		return
	}

	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.AlbumCreate(user))
	c.Render(http.StatusOK, rend)
}

func (r *ViewRouter) AlbumsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	user := r.auth.UserFromContext(ctx)

	var albms []*album.Album

	if !user.HasMemberRights() {
		albms = make([]*album.Album, 0)
	} else {
		var err error
		albms, err = r.albm.ListAlbums(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.Albums(user, albms))
	c.Render(http.StatusOK, rend)
}

func (r *ViewRouter) AlbumViewHandler(c *gin.Context) {
	ctx := c.Request.Context()
	user := r.auth.UserFromContext(ctx)

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id should not be empty"})
		return
	}

	almbID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	albm, err := r.albm.GetAlbum(ctx, almbID)
	if errors.Is(err, auth.ErrNotAuthorized) || errors.Is(err, auth.ErrNoMemberRights) || errors.Is(err, albumservice.ErrNotOwner) {
		c.Redirect(http.StatusFound, "/view/albums")
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	srch := search.NewPlantSearch()

	albumFilter := search.NewPlantAlbumFilter(albm.ID(), nil)
	srch.AddFilter(albumFilter)

	plants, err := r.srch.SearchPlants(ctx, srch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	plantMap := make(map[uuid.UUID]*searchservice.SearchPlant)
	for _, plnt := range plants {
		plantMap[plnt.ID] = plnt
		plnt.MainPhoto.URL = r.plantMedia.GetUrl(plnt.MainPhoto.URL)
	}
	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.AlbumView(user, albm, plantMap))
	c.Render(http.StatusOK, rend)
}

func (r *ViewRouter) AlbumUpdateHandler(c *gin.Context) {
	ctx := c.Request.Context()
	user := r.auth.UserFromContext(ctx)

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id should not be empty"})
		return
	}

	albmID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	albm, err := r.albm.GetAlbum(ctx, albmID)
	if errors.Is(err, auth.ErrNotAuthorized) || errors.Is(err, auth.ErrNoMemberRights) || errors.Is(err, albumservice.ErrNotOwner) {
		c.Redirect(http.StatusFound, "/view/albums")
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.AlbumUpdate(user, albm))
	c.Render(http.StatusOK, rend)
}
