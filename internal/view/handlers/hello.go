package handlers

import (
	"PlantSite/internal/view/components"
	"PlantSite/internal/view/gintemplrenderer"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HelloHandler(c *gin.Context) {
	r := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.Hello("Gin"))
	c.Render(http.StatusOK, r)
}
