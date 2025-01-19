package pdcp_service

import (
    "net/http"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    forwardRequests = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "pdcp_forward_requests_total",
            Help: "Number of data forwarding requests handled by PDCP",
        },
        []string{"slice_id"},
    )
)

func init() {
    prometheus.MustRegister(forwardRequests)
}

func incForwardCounter(sliceID string) {
    forwardRequests.WithLabelValues(sliceID).Inc()
}

func MetricsHandler() http.Handler {
    return promhttp.Handler()
}
