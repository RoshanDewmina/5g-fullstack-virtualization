package main

import (
    "log"
    "net/http"
    "os"

    "github.com/youruser/5g-fullstack-virtualization/slice-manager"
)

func main() {
    port := getPort()

    http.HandleFunc("/health", slice_manager.HealthHandler)
    http.HandleFunc("/createSlice", slice_manager.CreateSliceHandler)
    http.HandleFunc("/deleteSlice", slice_manager.DeleteSliceHandler)
    http.Handle("/metrics", slice_manager.MetricsHandler())

    log.Printf("Starting Slice Manager on port %s ...", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal("Slice Manager failed:", err)
    }
}

func getPort() string {
    port := os.Getenv("SLICE_MANAGER_PORT")
    if port == "" {
        port = "8084"
    }
    return port
}
