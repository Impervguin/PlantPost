package view

import (
	"PlantSite/internal/view/components"
	"PlantSite/internal/view/gintemplrenderer"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *ViewRouter) LoginHandler(c *gin.Context) {
	ctx := c.Request.Context()
	user := r.auth.UserFromContext(ctx)
	if user.IsAuthenticated() {
		c.Redirect(http.StatusFound, "/view")
		return
	}

	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.Login())
	c.Render(http.StatusOK, rend)
}

func (r *ViewRouter) LogoutHandler(c *gin.Context) {
	ctx := c.Request.Context()
	user := r.auth.UserFromContext(ctx)
	if !user.IsAuthenticated() {
		c.Redirect(http.StatusFound, "/view")
		return
	}

	err := r.auth.Logout(ctx)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/view/login")
}
