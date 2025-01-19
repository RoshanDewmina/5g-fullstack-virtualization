package rrm_service

import (
    "encoding/json"
    "net/http"
    "sync"
    "time"
)

// We’ll keep an in-memory map of slice allocations
// Key: sliceID, Value: allocated bandwidth (MHz, or any resource unit)
var (
    sliceAllocations = make(map[string]string)
    mutex            sync.Mutex
)

// AllocateHandler simulates resource allocation
// Request example: { "sliceID": "slice1", "requiredBW": "10MHz" }
type AllocationRequest struct {
    SliceID    string `json:"sliceID"`
    RequiredBW string `json:"requiredBW"`
    Priority   int    `json:"priority"`
}

type AllocationResponse struct {
    Status             string `json:"status"`
    ResourcesAllocated string `json:"resourcesAllocated"`
    Timestamp          string `json:"timestamp"`
    ErrorMsg           string `json:"errorMsg,omitempty"`
}

func AllocateHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var req AllocationRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
        return
    }

    // For demonstration, "allocate" means we just store the value in memory
    mutex.Lock()
    sliceAllocations[req.SliceID] = req.RequiredBW
    allocated := req.RequiredBW
    mutex.Unlock()

    resp := AllocationResponse{
        Status:             "OK",
        ResourcesAllocated: allocated,
        Timestamp:          time.Now().Format(time.RFC3339),
    }

    incAllocationCounter(req.SliceID) // Prometheus counter
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("RRM is healthy"))
}
