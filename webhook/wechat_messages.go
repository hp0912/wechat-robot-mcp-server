package webhook

import (
	"encoding/json"
	"io"
	"net/http"

	"wechat-robot-mcp-server/protobuf"
)

type WeChatMessageResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func OnWeChatMessages(w http.ResponseWriter, r *http.Request) {
	var req protobuf.WeChatMessage

	// 只接受 POST 请求
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(WeChatMessageResponse{
			Code:    http.StatusMethodNotAllowed,
			Message: "method not allowed, only POST is supported",
		})
		return
	}

	// 读取并解析请求体中的 JSON 为结构化的 WeChatMessage
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(WeChatMessageResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to read request body",
		})
		return
	}
	defer r.Body.Close()

	if len(body) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(WeChatMessageResponse{
			Code:    http.StatusBadRequest,
			Message: "empty request body",
		})
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(WeChatMessageResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid JSON body",
		})
		return
	}

	// TODO: 在这里继续处理解析后的 req（如入库、业务逻辑等）

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(WeChatMessageResponse{
		Code:    http.StatusOK,
		Message: "ok",
	})
}
