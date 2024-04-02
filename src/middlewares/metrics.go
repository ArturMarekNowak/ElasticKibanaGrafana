package middlewares

import (
	"ElasticKibanaGrafana/src/models"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

func Metrics(m *models.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		m.Requests.With(prometheus.Labels{"method": c.Request.Method, "path": c.Request.RequestURI}).Inc()
		start := time.Now()
		c.Next()
		stop := time.Since(start).Seconds()
		statusCode := strconv.Itoa(c.Writer.Status())
		m.Duration.With(prometheus.Labels{"method": c.Request.Method, "status": statusCode}).Observe(stop)
	}
}
