package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"wechat-robot-mcp-server/config"
	"wechat-robot-mcp-server/middleware"
	"wechat-robot-mcp-server/tools"
	"wechat-robot-mcp-server/webhook"
)

var Version = "unknown"

func main() {
	log.Printf("[MCP Server]启动 版本: %s", Version)

	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "wechat-robot-mcp-server",
		Version: "1.0.0",
	}, nil)
	server.AddReceivingMiddleware(middleware.TenantMiddleware)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "ChatRoomSummary",
		Description: "微信群聊总结，当用户想总结群聊内容时，可以调用该工具。",
	}, tools.ChatRoomSummary)

	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{})

	mux := http.NewServeMux()
	mux.Handle("/api/v1/messages", http.HandlerFunc(webhook.OnWeChatMessages))
	mux.Handle("/mcp", handler)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.MCPServerPort), mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
