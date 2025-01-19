package mac_service

import (
    "encoding/json"
    "log"
    "net/http"
    "strings"
    "time"
)

// Mock advanced scheduling logic:
// 1. Query RRM for the allocated resources for a slice
// 2. Perform a “scheduling decision” based on priority or existing load
// 3. Return scheduling info
type ScheduleRequest struct {
    SliceID string `json:"sliceID"`
}

type ScheduleResponse struct {
    Status       string `json:"status"`
    ScheduledBW  string `json:"scheduledBW"`
    SchedulingTs string `json:"schedulingTs"`
}

func ScheduleHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var req ScheduleRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
        return
    }

    // “Get allocated resources” from RRM
    // This is a naive example: we simply GET from RRM with the slice ID
    rrmResp, err := http.Get("http://rrm-service:8081/allocate?sliceID=" + req.SliceID)
    if err != nil {
        log.Println("Failed to contact RRM for scheduling:", err)
        http.Error(w, "RRM contact failed", http.StatusInternalServerError)
        return
    }
    // We’re calling a POST endpoint with GET, which is obviously not correct for real logic.
    // For demonstration, maybe you’d store slice allocations in RRM and expose a GET.

    // In a real system, parse actual JSON. Here, we skip detailed logic for brevity.

    // “Schedule” logic: assume we can use the same allocated resources
    // In real logic, we might have a fairness or priority-based algorithm
    // that changes how resources are split among multiple UEs/slices.

    response := ScheduleResponse{
        Status:       "Scheduled",
        ScheduledBW:  "Using allocated resources from RRM",
        SchedulingTs: time.Now().Format(time.RFC3339),
    }

    incScheduleCounter(req.SliceID) // Prometheus
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("MAC is healthy"))
}
