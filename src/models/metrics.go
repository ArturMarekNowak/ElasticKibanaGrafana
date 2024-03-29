package models

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	Requests *prometheus.CounterVec
	Duration *prometheus.HistogramVec
}
