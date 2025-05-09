package view

import (
	authservice "PlantSite/internal/services/auth-service"
	searchservice "PlantSite/internal/services/search-service"
	"PlantSite/internal/view/components"
	"PlantSite/internal/view/gintemplrenderer"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MediaUrlStrategy interface {
	GetUrl(path string) string
}

type ViewRouter struct {
	StaticPath string
	auth       *authservice.AuthService
	srch       *searchservice.SearchService
	plantMedia MediaUrlStrategy
	postMedia  MediaUrlStrategy
}

func (r *ViewRouter) Init(router *gin.RouterGroup, staticPath string, auth *authservice.AuthService, srch *searchservice.SearchService, plantMedia MediaUrlStrategy, postMedia MediaUrlStrategy) {
	r.auth = auth
	r.srch = srch
	r.plantMedia = plantMedia
	r.postMedia = postMedia
	r.StaticPath = staticPath
	router.StaticFS("/static", http.Dir(staticPath))

	gr := router.Group("/view")
	gr.GET("/", r.IndexHandler)
	gr.GET("/login", r.LoginHandler)
	gr.GET("/register", r.RegisterHandler)
	gr.GET("/logout", r.LogoutHandler)
	gr.GET("/plants", r.PlantsHandler)
}

func (r *ViewRouter) IndexHandler(c *gin.Context) {
	ctx := c.Request.Context()
	user := r.auth.UserFromContext(ctx)
	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.Index(user))
	c.Render(http.StatusOK, rend)
}
