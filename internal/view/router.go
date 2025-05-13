package view

import (
	albumservice "PlantSite/internal/services/album-service"
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
	albm       *albumservice.AlbumService
	plantMedia MediaUrlStrategy
	postMedia  MediaUrlStrategy
}

func (r *ViewRouter) Init(
	router *gin.RouterGroup,
	staticPath string,
	auth *authservice.AuthService,
	srch *searchservice.SearchService,
	albm *albumservice.AlbumService,
	plantMedia MediaUrlStrategy,
	postMedia MediaUrlStrategy) {
	r.auth = auth
	r.srch = srch
	r.albm = albm
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
	gr.GET("/plant/create", r.CreatePlantHandler)
	gr.GET("/plant/:id", r.PlantViewHandler)
	gr.GET("/plant/:id/update", r.UpdatePlantHandler)

	gr.GET("/posts", r.PostsHandler)
	gr.GET("/post/:id", r.PostViewHandler)
	gr.GET("/post/create", r.CreatePostHandler)
	gr.GET("/post/:id/update", r.UpdatePostHandler)

	gr.GET("/albums", r.AlbumsHandler)
	gr.GET("/album/:id", r.AlbumViewHandler)
	gr.GET("/album/create", r.AlbumsCreateHandler)
	gr.GET("/album/:id/update", r.AlbumUpdateHandler)
}

func (r *ViewRouter) IndexHandler(c *gin.Context) {
	ctx := c.Request.Context()
	user := r.auth.UserFromContext(ctx)
	rend := gintemplrenderer.New(c.Request.Context(), http.StatusOK, components.Index(user))
	c.Render(http.StatusOK, rend)
}
