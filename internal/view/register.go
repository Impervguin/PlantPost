package view

import (
	"PlantSite/internal/view/components"
	"PlantSite/internal/view/gintemplrenderer"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *ViewRouter) RegisterHandler(c *gin.Context) {
	ctx := c.Request.Context()
	user := r.auth.UserFromContext(ctx)
	if user.IsAuthenticated() {
		c.Redirect(http.StatusFound, "/view")
		return
	}

	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.Register())
	c.Render(http.StatusOK, rend)
}
