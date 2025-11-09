package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var Version = "unknown"

func main() {
	log.Printf("[MCP Server]启动 版本: %s", Version)

	// 加载配置
	if err := LoadConfig(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "wechat-robot-mcp-server",
		Version: "1.0.0",
	}, nil)
	server.AddReceivingMiddleware(TenantMiddleware)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "ChatRoomSummary",
		Description: "微信群聊总结，当用户想总结群聊内容时，可以调用该工具。",
	}, ChatRoomSummary)

	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{})

	mux := http.NewServeMux()
	mux.Handle("/api/v1/messages", http.HandlerFunc(onWeChatMessages))
	mux.Handle("/mcp", handler)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", MCPServerPort), mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
