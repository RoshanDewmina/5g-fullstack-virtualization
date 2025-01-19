package main

import (
    "log"
    "net/http"
    "os"

    "github.com/youruser/5g-fullstack-virtualization/pdcp-service"
)

func main() {
    port := getPort()

    // Register endpoints
    http.HandleFunc("/health", pdcp_service.HealthHandler)
    http.HandleFunc("/forward", pdcp_service.ForwardDataHandler)
    http.Handle("/metrics", pdcp_service.MetricsHandler())

    log.Printf("Starting PDCP Service on port %s ...", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal("PDCP Service failed:", err)
    }
}

func getPort() string {
    port := os.Getenv("PDCP_SERVICE_PORT")
    if port == "" {
        port = "8083"
    }
    return port
}
