package mac_service

import (
    "net/http"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    scheduleRequests = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "mac_schedule_requests_total",
            Help: "Number of scheduling requests handled by MAC",
        },
        []string{"slice_id"},
    )
)

func init() {
    prometheus.MustRegister(scheduleRequests)
}

func incScheduleCounter(sliceID string) {
    scheduleRequests.WithLabelValues(sliceID).Inc()
}

func MetricsHandler() http.Handler {
    return promhttp.Handler()
}
