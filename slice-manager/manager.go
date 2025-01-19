package slice_manager

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "sync"
    "time"
)

// Example in-memory store of slices.
var (
    slices = make(map[string]SliceConfig)
    mu     sync.Mutex
)

type SliceConfig struct {
    SliceID    string `json:"sliceID"`
    RequiredBW string `json:"requiredBW"`
    Priority   int    `json:"priority"`
}

// CreateSliceHandler calls RRM to allocate resources, then sets up MAC, PDCP, etc.
func CreateSliceHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    defer r.Body.Close()
    var sc SliceConfig
    if err := json.NewDecoder(r.Body).Decode(&sc); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    mu.Lock()
    if _, exists := slices[sc.SliceID]; exists {
        mu.Unlock()
        http.Error(w, "Slice already exists", http.StatusConflict)
        return
    }
    mu.Unlock()

    // 1. Call RRM /allocate
    allocReq := map[string]interface{}{
        "sliceID":    sc.SliceID,
        "requiredBW": sc.RequiredBW,
        "priority":   sc.Priority,
    }
    err := postJSON("http://rrm-service:8081/allocate", allocReq)
    if err != nil {
        log.Println("Error allocating resources in RRM:", err)
        http.Error(w, "Failed RRM allocation", http.StatusInternalServerError)
        return
    }

    // 2. Call MAC /schedule (optionally with the slice ID).
    schedReq := map[string]interface{}{
        "sliceID": sc.SliceID,
    }
    err = postJSON("http://mac-service:8082/schedule", schedReq)
    if err != nil {
        log.Println("Error scheduling in MAC:", err)
        http.Error(w, "Failed MAC scheduling", http.StatusInternalServerError)
        return
    }

    // 3. Possibly call PDCP /forward if we want to initialize something
    // For demonstration, do nothing special here. Or do a small call:
    forwardReq := map[string]interface{}{
        "sliceID":   sc.SliceID,
        "plainData": "INITIAL DATA",
    }
    err = postJSON("http://pdcp-service:8083/forward", forwardReq)
    if err != nil {
        log.Println("Error calling PDCP:", err)
        http.Error(w, "Failed PDCP init", http.StatusInternalServerError)
        return
    }

    mu.Lock()
    slices[sc.SliceID] = sc
    mu.Unlock()

    incSliceCreateCounter(sc.SliceID)
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte(fmt.Sprintf("Slice %s created successfully", sc.SliceID)))
}

func DeleteSliceHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }
    defer r.Body.Close()

    var req struct {
        SliceID string `json:"sliceID"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    mu.Lock()
    if _, exists := slices[req.SliceID]; !exists {
        mu.Unlock()
        http.Error(w, "Slice not found", http.StatusNotFound)
        return
    }
    delete(slices, req.SliceID)
    mu.Unlock()

    incSliceDeleteCounter(req.SliceID)
    // No real “deallocation” calls here, but you could call RRM/MAC/PDCP cleanup endpoints
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(fmt.Sprintf("Slice %s deleted successfully", req.SliceID)))
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Slice Manager is healthy"))
}

// Helper to make a POST with JSON body
func postJSON(url string, data interface{}) error {
    body, err := json.Marshal(data)
    if err != nil {
        return err
    }

    resp, err := http.Post(url, "application/json", bytes.NewReader(body))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 300 {
        b, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("status %d: %s", resp.StatusCode, string(b))
    }
    return nil
}
