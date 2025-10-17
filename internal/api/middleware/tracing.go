package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	meter          = otel.Meter("gin-server")
	requestCounter metric.Int64Counter
)

func init() {
	var err error
	requestCounter, err = meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total HTTP requests"),
	)
	if err != nil {
		panic(err)
	}
}

func OpenTelemetryMiddleware() gin.HandlerFunc {
	// Базовый middleware от OpenTelemetry
	baseMiddleware := otelgin.Middleware("plant-site-api")

	return func(c *gin.Context) {
		start := time.Now()

		// Вызываем базовый middleware
		baseMiddleware(c)

		// Замер времени выполнения
		c.Next()

		duration := time.Since(start).Milliseconds()

		// Записываем метрику времени выполнения
		// (можно добавить гистограмму для latency)
		requestCounter.Add(c.Request.Context(), 1, metric.WithAttributes(
			attribute.String("method", c.Request.Method),
			attribute.String("path", c.Request.URL.Path),
			attribute.Int64("duration", duration),
		))
	}
}
