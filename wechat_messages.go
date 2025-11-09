package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type PostMessageRequest struct {
	Text string `json:"text"`
}

type PostMessageResponse struct {
	Success bool   `json:"success"`
	Echo    string `json:"echo,omitempty"`
	Error   string `json:"error,omitempty"`
}

func onWeChatMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(PostMessageResponse{
			Success: false,
			Error:   "method not allowed, use POST",
		})
		return
	}

	defer r.Body.Close()

	var req PostMessageRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(PostMessageResponse{
			Success: false,
			Error:   "invalid JSON payload",
		})
		return
	}

	if req.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(PostMessageResponse{
			Success: false,
			Error:   "text is required",
		})
		return
	}

	log.Printf("收到消息: %q", req.Text)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(PostMessageResponse{
		Success: true,
		Echo:    req.Text,
	})
}
