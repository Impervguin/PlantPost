package middleware

import (
	authapi "PlantSite/internal/api/auth-api"
	"PlantSite/internal/models/auth"
	authservice "PlantSite/internal/services/auth-service"
	"PlantSite/internal/utils/logs"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware(s authservice.AuthServiceContract) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sessID uuid.UUID = uuid.Nil
		if cookie, err := c.Request.Cookie(authapi.SessionCookieName); err == nil {
			sessID, err = uuid.Parse(cookie.Value)
			if err != nil {
				sessID = uuid.Nil
			}
		}
		ctx := s.Authenticate(c.Request.Context(), sessID)
		c.Request = c.Request.WithContext(ctx)

		user := s.UserFromContext(ctx)

		if _, ok := user.(*auth.NoAuthUser); ok {
			logs.Infow("request unauthenticated", "request_id", c.GetString(RequestIDKey))
		} else {
			logs.Infow("request authenticated", "request_id", c.GetString(RequestIDKey), "user_id", user.ID())
		}

		if _, ok := user.(*auth.NoAuthUser); ok {
			c.SetCookie(authapi.SessionCookieName, "", -1, "/", "", false, false)
		}

		c.Next()
	}
}
