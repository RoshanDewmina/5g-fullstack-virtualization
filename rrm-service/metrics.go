package rrm_service

import (
    "net/http"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    allocationRequests = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "rrm_allocation_requests_total",
            Help: "Number of resource allocation requests handled by RRM",
        },
        []string{"slice_id"},
    )
)

func init() {
    prometheus.MustRegister(allocationRequests)
}

func incAllocationCounter(sliceID string) {
    allocationRequests.WithLabelValues(sliceID).Inc()
}

func MetricsHandler() http.Handler {
    return promhttp.Handler()
}
