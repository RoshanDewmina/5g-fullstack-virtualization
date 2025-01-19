package main

import (
    "log"
    "net/http"
    "os"

    "github.com/youruser/5g-fullstack-virtualization/rrm-service"
)

func main() {
    port := getPort()

    // Register endpoints
    http.HandleFunc("/health", rrm_service.HealthHandler)
    http.HandleFunc("/allocate", rrm_service.AllocateHandler)
    http.Handle("/metrics", rrm_service.MetricsHandler())

    log.Printf("Starting RRM Service on port %s ...", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal("RRM Service failed:", err)
    }
}

func getPort() string {
    port := os.Getenv("RRM_SERVICE_PORT")
    if port == "" {
        port = "8081"
    }
    return port
}
