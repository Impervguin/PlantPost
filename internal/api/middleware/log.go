package middleware

import (
	"PlantSite/internal/utils/logs"

	"github.com/gin-gonic/gin"
)

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		reqID := c.GetString(RequestIDKey)
		if reqID == "" {
			reqID = "unknown"
		}
		logs.Infow("request started", "request_id", reqID, "method", c.Request.Method, "path", c.Request.URL.Path)
		c.Next()
		errs := c.Errors

		if len(errs) > 0 {
			logs.Errorw("request failed", "request_id", reqID, "method", c.Request.Method, "path", c.Request.URL.Path, "error", errs)
		}
		logs.Infow("request finished", "request_id", reqID, "method", c.Request.Method, "path", c.Request.URL.Path)
	}
}
