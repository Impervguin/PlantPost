package metricsapi

import (
	"PlantSite/internal/infra/metrics"

	"github.com/gin-gonic/gin"
)

type MetricsRouter struct{}

func (r *MetricsRouter) Init(router *gin.RouterGroup) {
	gr := router.Group("/metrics")
	gr.GET("/get", r.Get)
}

func (r *MetricsRouter) Get(c *gin.Context) {
	c.JSON(200, gin.H{"metrics": metrics.GetMetrics()})
}
