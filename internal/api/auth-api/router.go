package authapi

import (
	authservice "PlantSite/internal/services/auth-service"

	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	auth *authservice.AuthService
}

func (r *AuthRouter) Init(router *gin.RouterGroup, auth *authservice.AuthService) {
	r.auth = auth
	gr := router.Group("/auth")
	gr.POST("/login", r.Login)
}

func (r *AuthRouter) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	sessID, err := r.auth.Login(ctx, req.Username, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie(SessionCookieName, sessID.String(), 0, "/", "", false, true)
	c.JSON(200, gin.H{})
}
