package view

import (
	plantsquery "PlantSite/internal/api-utils/query-filters/plants-query"
	"PlantSite/internal/view/components"
	"PlantSite/internal/view/gintemplrenderer"
	"net/http"

	"github.com/gin-gonic/gin"
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
