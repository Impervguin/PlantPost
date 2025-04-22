package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"
const RequestIDKey = "RequestID"

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set(RequestIDKey, requestID)

		c.Writer.Header().Set(RequestIDHeader, requestID)

		c.Next()
	}
}
