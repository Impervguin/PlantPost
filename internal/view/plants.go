package view

import (
	plantsquery "PlantSite/internal/api-utils/query-filters/plants-query"
	"PlantSite/internal/view/components"
	"PlantSite/internal/view/gintemplrenderer"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *ViewRouter) PlantsHandler(c *gin.Context) {
	srch, err := plantsquery.ParseGinQueryPlantSearch(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	plnts, err := r.srch.SearchPlants(c.Request.Context(), srch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, plnt := range plnts {
		plnt.MainPhoto.URL = r.plantMedia.GetUrl(plnt.MainPhoto.URL)
	}

	ctx := c.Request.Context()
	user := r.auth.UserFromContext(ctx)

	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.Plants(user, plnts))
	c.Render(http.StatusOK, rend)
}

func (r *ViewRouter) CreatePlantHandler(c *gin.Context) {
	ctx := c.Request.Context()
	user := r.auth.UserFromContext(ctx)
	if !user.HasAuthorRights() {
		c.Redirect(http.StatusFound, "/view/plants")
		return
	}

	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.CreatePlant(user))
	c.Render(http.StatusOK, rend)
}

type plantView struct {
	ID string `uri:"id" binding:"required"`
}

func (r *ViewRouter) PlantViewHandler(c *gin.Context) {
	ctx := c.Request.Context()
	user := r.auth.UserFromContext(ctx)

	var plntView plantView

	if err := c.ShouldBindUri(&plntView); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := uuid.Parse(plntView.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plnt, err := r.srch.GetPlantByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	plnt.MainPhoto.URL = r.plantMedia.GetUrl(plnt.MainPhoto.URL)

	for _, photo := range plnt.Photos {
		photo.File.URL = r.plantMedia.GetUrl(photo.File.URL)
	}

	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.PlantView(user, plnt))
	c.Render(http.StatusOK, rend)
}

func (r *ViewRouter) UpdatePlantHandler(c *gin.Context) {
	ctx := c.Request.Context()
	user := r.auth.UserFromContext(ctx)
	if !user.HasAuthorRights() {
		c.Redirect(http.StatusFound, "/view/plants")
		return
	}

	var plntView plantView

	if err := c.ShouldBindUri(&plntView); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := uuid.Parse(plntView.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plnt, err := r.srch.GetPlantByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	plnt.MainPhoto.URL = r.plantMedia.GetUrl(plnt.MainPhoto.URL)

	for _, photo := range plnt.Photos {
		photo.File.URL = r.plantMedia.GetUrl(photo.File.URL)
	}

	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.UpdatePlantSpecification(user, plnt))
	c.Render(http.StatusOK, rend)
}
