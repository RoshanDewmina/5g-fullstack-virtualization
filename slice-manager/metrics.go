package slice_manager

import (
    "net/http"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    sliceCreateCount = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "slice_manager_create_requests_total",
            Help: "Number of slice creation requests handled",
        },
        []string{"slice_id"},
    )

    sliceDeleteCount = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "slice_manager_delete_requests_total",
            Help: "Number of slice deletion requests handled",
        },
        []string{"slice_id"},
    )
)

func init() {
    prometheus.MustRegister(sliceCreateCount, sliceDeleteCount)
}

func incSliceCreateCounter(sliceID string) {
    sliceCreateCount.WithLabelValues(sliceID).Inc()
}

func incSliceDeleteCounter(sliceID string) {
    sliceDeleteCount.WithLabelValues(sliceID).Inc()
}

func MetricsHandler() http.Handler {
    return promhttp.Handler()
}
