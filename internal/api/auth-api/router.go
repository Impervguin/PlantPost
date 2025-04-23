package authapi

import (
	authservice "PlantSite/internal/services/auth-service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	auth *authservice.AuthService
}

func (r *AuthRouter) Init(router *gin.RouterGroup, auth *authservice.AuthService) {
	r.auth = auth
	gr := router.Group("/auth")
	gr.POST("/login", r.Login)
	gr.POST("/register", r.Register)
	gr.POST("/logout", r.Logout)
}

// Login Handler
// @Summary User login
// @Description Authenticates a user and creates a session
// @Tags auth
// @Accept json
// @Accept mpfd
// @Param request body LoginRequest true "Login credentials"
// @Success 200 "Session for user created"
// @Failure 400 "Wrong input parameters"
// @Failure 401 "Auth error"
// @Router /auth/login [post]
func (r *AuthRouter) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sessID, err := r.auth.Login(ctx, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie(SessionCookieName, sessID.String(), 0, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
}

// Register Handler
// @Summary Register a new user
// @Description Registers a new user. Requires login afterwards
// @Tags auth
// @Accept json
// @Accept mpfd
// @Param request body RegisterRequest true "Register credentials"
// @Success 200 "User registered"
// @Failure 400 "Wrong input parameters"
// @Failure 401 "Auth error"
// @Router /auth/register [post]
func (r *AuthRouter) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := r.auth.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// Logouts Handler
// @Summary Logouts a user
// @Description Logouts a user
// @Tags auth
// @Success 200 "User logged out"
// @Failure 401 "Logout error"
// @Router /auth/logout [post]
func (r *AuthRouter) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	err := r.auth.Logout(ctx)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
