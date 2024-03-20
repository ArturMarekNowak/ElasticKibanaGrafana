package main

import (
	"ElasticKibanaGrafanaJaeger/src/controllers"
	"ElasticKibanaGrafanaJaeger/src/middlewares"
	"github.com/gin-gonic/gin"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"log"
	"os"
)

func ConfigureServer() {

	router := gin.Default()
	logger := ConfigureLogging()
	ConfigureMiddlewares(router, logger)
	ConfigureEndpoints(router)
	err := router.Run(":8080")
	if err != nil {
		log.Fatal("Couldn't start host")
	}
}

func ConfigureMiddlewares(router *gin.Engine, logger *zap.Logger) {
	router.Use(middlewares.Logging(logger))
	router.Use(middlewares.Metrics)
	router.Use(middlewares.Tracing)
}

func ConfigureEndpoints(router *gin.Engine) {
	router.GET("/helloWorld", controllers.GetHelloWorld)
}

func ConfigureLogging() *zap.Logger {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	return logger
}

func main() {
	ConfigureServer()
}
