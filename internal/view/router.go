package view

import (
	authservice "PlantSite/internal/services/auth-service"
	"PlantSite/internal/view/components"
	"PlantSite/internal/view/gintemplrenderer"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ViewRouter struct {
	StaticPath string
	auth       *authservice.AuthService
}

func (r *ViewRouter) Init(router *gin.RouterGroup, staticPath string, auth *authservice.AuthService) {
	r.auth = auth
	r.StaticPath = staticPath
	router.StaticFS("/static", http.Dir(staticPath))

	gr := router.Group("/view")
	gr.GET("/", r.IndexHandler)
	gr.GET("/login", r.LoginHandler)
	gr.GET("/register", r.RegisterHandler)
	gr.GET("/logout", r.LogoutHandler)
}

func (r *ViewRouter) IndexHandler(c *gin.Context) {
	ctx := c.Request.Context()
	user := r.auth.UserFromContext(ctx)
	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.Index(user))
	c.Render(http.StatusOK, rend)
}
