package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

func Tracing(c *gin.Context) {
	_, span := otel.Tracer(c.Request.Method+c.Request.RequestURI).Start(c, "request")
	defer span.End()
	c.Next()
}
