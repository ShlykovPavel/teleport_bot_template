package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics структура для хранения метрик Prometheus
type Metrics struct {
	HttpRequestsTotal   *prometheus.CounterVec
	HttpRequestDuration *prometheus.HistogramVec
}

// NewMetrics создает новый набор метрик
func NewMetrics() *Metrics {
	return &Metrics{
		HttpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"}),
		HttpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests",
				Buckets: []float64{0.01, 0.1, 0.5, 1, 2, 5},
			},
			[]string{"method", "path"},
		),
	}
}

// InitMetrics инициализирует и регистрирует метрики в Prometheus
func InitMetrics() *Metrics {
	metrics := NewMetrics()
	prometheus.MustRegister(
		metrics.HttpRequestsTotal,
		metrics.HttpRequestDuration,
	)
	return metrics
}
