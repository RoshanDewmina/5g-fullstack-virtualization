package main

import (
    "log"
    "net/http"
    "os"

    "github.com/youruser/5g-fullstack-virtualization/mac-service"
)

func main() {
    port := getPort()

    // Register endpoints
    http.HandleFunc("/health", mac_service.HealthHandler)
    http.HandleFunc("/schedule", mac_service.ScheduleHandler)
    http.Handle("/metrics", mac_service.MetricsHandler())

    log.Printf("Starting MAC Service on port %s ...", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal("MAC Service failed:", err)
    }
}

func getPort() string {
    port := os.Getenv("MAC_SERVICE_PORT")
    if port == "" {
        port = "8082"
    }
    return port
}
