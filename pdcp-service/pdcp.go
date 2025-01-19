package pdcp_service

import (
    "encoding/json"
    "net/http"
    "time"
)

// This simulates PDCP’s data forwarding with a minimal ciphering placeholder.
type ForwardRequest struct {
    SliceID   string `json:"sliceID"`
    PlainData string `json:"plainData"`
}

type ForwardResponse struct {
    Status     string `json:"status"`
    CipherData string `json:"cipherData"`
    Timestamp  string `json:"timestamp"`
}

func ForwardDataHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var req ForwardRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Mock encryption: reverse the string (just for demonstration)
    cipherData := reverseString(req.PlainData)

    resp := ForwardResponse{
        Status:     "Data forwarded",
        CipherData: cipherData,
        Timestamp:  time.Now().Format(time.RFC3339),
    }

    incForwardCounter(req.SliceID)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func reverseString(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("PDCP is healthy"))
}
