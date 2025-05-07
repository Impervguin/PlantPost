package view

import (
	authservice "PlantSite/internal/services/auth-service"
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
	gr.GET("/login", r.LoginHandler)
	gr.GET("/register", r.RegisterHandler)
}
