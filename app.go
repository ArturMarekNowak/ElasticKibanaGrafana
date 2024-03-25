package main

import (
	"ElasticKibanaGrafanaJaeger/src/controllers"
	"ElasticKibanaGrafanaJaeger/src/middlewares"
	"ElasticKibanaGrafanaJaeger/src/models"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.elastic.co/ecszap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func ConfigureServer() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	reg := prometheus.NewRegistry()
	router := gin.Default()
	logger := ConfigureLogging()
	metrics := ConfigureMetrics(reg)
	otelShutdown, err := setupOTelSDK(ctx)
	if err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()
	ConfigureMiddlewares(router, logger, metrics)
	ConfigureEndpoints(router, reg)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	if err != nil {
		log.Fatal("Couldn't start host")
	}
	srvErr := make(chan error, 1)

	select {
	case err = <-srvErr:
		return
	case <-ctx.Done():
		stop()
	}

	err = srv.Shutdown(context.Background())
}

func setupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)
	tracerProvider, err := newTraceProvider(ctx)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)
	return
}

func newTraceProvider(ctx context.Context) (*trace.TracerProvider, error) {
	traceExporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithInsecure())
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
	)
	return traceProvider, nil
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
	router.Use(middlewares.Tracing)
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
