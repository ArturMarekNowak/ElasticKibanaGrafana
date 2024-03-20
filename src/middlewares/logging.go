package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logging(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.With(zap.String("method", c.Request.Method)).With(zap.String("path", c.Request.RequestURI)).Info("request processing began")
		c.Next()
		logger.With(zap.String("method", c.Request.Method)).With(zap.String("path", c.Request.RequestURI)).Info("request processed")
	}
}
