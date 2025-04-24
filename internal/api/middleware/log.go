package middleware

import (
	"github.com/gin-gonic/gin"
)

type MiddlewareLogger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
}

const LoggerKey = "logger"

func LogMiddleware(l MiddlewareLogger) gin.HandlerFunc {
	return func(c *gin.Context) {

		reqID := c.GetString(RequestIDKey)
		if reqID == "" {
			reqID = "unknown"
		}
		l.Infow("request started", "request_id", reqID, "method", c.Request.Method, "path", c.Request.URL.Path)
		c.Set(LoggerKey, l)
		c.Next()
		errs := c.Errors

		if len(errs) > 0 {
			l.Errorw("request failed", "request_id", reqID, "method", c.Request.Method, "path", c.Request.URL.Path, "error", errs)
		}
		l.Infow("request finished", "request_id", reqID, "method", c.Request.Method, "path", c.Request.URL.Path)
	}
}
