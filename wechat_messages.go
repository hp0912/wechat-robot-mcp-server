package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WeChatMessageResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func onWeChatMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(WeChatMessageResponse{
			Code:    500,
			Message: "method not allowed, use POST",
		})
		return
	}

	defer r.Body.Close()

	var req WeChatMessage
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(WeChatMessageResponse{
			Code:    500,
			Message: fmt.Sprintf("解析消息失败: %v", err),
		})
		return
	}

	_ = r.Header.Get("Robot-Code")
	// 处理消息逻辑 TODO:

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(WeChatMessageResponse{
		Code: 200,
	})
}
