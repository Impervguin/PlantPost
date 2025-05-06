package view

import (
	"PlantSite/internal/view/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ViewRouter struct {
	StaticPath string
}

func (r *ViewRouter) Init(router *gin.RouterGroup, staticPath string) {
	r.StaticPath = staticPath
	router.StaticFS("/static", http.Dir(staticPath))

	gr := router.Group("/view")
	gr.GET("/hello", handlers.HelloHandler)
}
