package api_gateway

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// In a real environment, you'd get these from config or environment
var (
    RRMServiceURL  = "http://rrm-service:8081"
    MACServiceURL  = "http://mac-service:8082"
    PDCPServiceURL = "http://pdcp-service:8083"
    SMServiceURL   = "http://slice-manager:8084"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("API Gateway is healthy"))
}

func CreateSliceHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Forward the request body to the Slice Manager
    // The Slice Manager orchestrates the creation logic
    defer r.Body.Close()
    body, err := readBody(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    resp, err := http.Post(SMServiceURL+"/createSlice", "application/json", bytes.NewReader(body))
    if err != nil {
        log.Println("Error calling Slice Manager:", err)
        http.Error(w, "Internal error contacting slice manager", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    w.WriteHeader(resp.StatusCode)
    w.Write([]byte("Slice creation invoked successfully"))
}

func DeleteSliceHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    defer r.Body.Close()
    body, err := readBody(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    req, err := http.NewRequest(http.MethodDelete, SMServiceURL+"/deleteSlice", bytes.NewReader(body))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Println("Error calling Slice Manager:", err)
        http.Error(w, "Internal error contacting slice manager", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    w.WriteHeader(resp.StatusCode)
    w.Write([]byte("Slice deletion invoked successfully"))
}

func readBody(r *http.Request) ([]byte, error) {
    buf := new(bytes.Buffer)
    _, err := buf.ReadFrom(r.Body)
    if err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

// Prometheus metrics endpoint
func MetricsHandler() http.Handler {
    return promhttp.Handler()
}
