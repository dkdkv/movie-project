// pkg/metrics/metrics.go
package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
)

var (
	MovieCreations = promauto.NewCounter(prometheus.CounterOpts{
		Name: "movie_creations_total",
		Help: "The total number of created movies",
	})

	MovieRetrievals = promauto.NewCounter(prometheus.CounterOpts{
		Name: "movie_retrievals_total",
		Help: "The total number of movie retrievals",
	})

	MovieUpdates = promauto.NewCounter(prometheus.CounterOpts{
		Name: "movie_updates_total",
		Help: "The total number of movie updates",
	})

	MovieDeletions = promauto.NewCounter(prometheus.CounterOpts{
		Name: "movie_deletions_total",
		Help: "The total number of movie deletions",
	})

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests in seconds",
		},
		[]string{"handler", "method"},
	)

	TotalRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"handler", "code", "method"},
	)
)

// InitMetrics initializes the metrics
func InitMetrics() {
	// Metrics are automatically registered via promauto
}

// UnaryServerInterceptor is a gRPC interceptor for Prometheus metrics
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start).Seconds()
	RequestDuration.WithLabelValues("grpc", info.FullMethod).Observe(duration)
	return resp, err
}
