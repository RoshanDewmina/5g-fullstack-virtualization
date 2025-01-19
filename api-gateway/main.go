package main

import (
    "log"
    "net/http"
    "os"

    "github.com/youruser/5g-fullstack-virtualization/api-gateway"
)

func main() {
    port := getPort()

    // Register HTTP handlers
    http.HandleFunc("/health", api_gateway.HealthHandler)
    // Protected endpoints
    // For example, we might protect the “/createSlice” or “/invoke” endpoints with JWT:
    http.HandleFunc("/createSlice", api_gateway.JWTAuth(api_gateway.CreateSliceHandler))
    http.HandleFunc("/deleteSlice", api_gateway.JWTAuth(api_gateway.DeleteSliceHandler))
    http.Handle("/metrics", api_gateway.MetricsHandler())

    log.Printf("Starting API Gateway on port %s ...", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal("API Gateway failed:", err)
    }
}

func getPort() string {
    port := os.Getenv("API_GATEWAY_PORT")
    if port == "" {
        port = "8080"
    }
    return port
}
