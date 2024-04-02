package main

import (
	"ElasticKibanaGrafana/src/controllers"
	"ElasticKibanaGrafana/src/middlewares"
	"ElasticKibanaGrafana/src/models"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func ConfigureServer() {
	reg := prometheus.NewRegistry()
	router := gin.Default()
	logger := ConfigureLogging()
	metrics := ConfigureMetrics(reg)
	ConfigureMiddlewares(router, logger, metrics)
	ConfigureEndpoints(router, reg)
	StartServer(router)
}

func StartServer(router *gin.Engine) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Couldn't start host: %s\n", err)
		}
	}()

	select {
	case <-ctx.Done():
		stop()
	}

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Couldn't start host:: %s\n", err)
	}
}

func ConfigureMetrics(reg prometheus.Registerer) *models.Metrics {
	m := &models.Metrics{
		Requests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "app",
			Name:      "http_requests_count",
			Help:      "Number of requests.",
		}, []string{"method", "path"}),
		Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "app",
			Name:      "http_requests_duration",
			Help:      "Duration of the request.",
			Buckets:   []float64{0.1, 0.15, 0.2, 0.25, 0.3},
		}, []string{"method", "status"}),
	}
	reg.MustRegister(m.Requests, m.Duration)
	return m
}

func ConfigureMiddlewares(router *gin.Engine, logger *zap.Logger, metrics *models.Metrics) {
	router.Use(middlewares.Logging(logger))
	router.Use(middlewares.Metrics(metrics))
}

func ConfigureEndpoints(router *gin.Engine, reg prometheus.Gatherer) {
	router.GET("/helloWorld", controllers.GetHelloWorld)
	router.GET("/metrics", gin.WrapH(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})))
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
