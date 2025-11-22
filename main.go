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
		Title:       "群聊总结",
		Description: "微信群聊总结，当用户想总结群聊内容时，可以调用该工具。",
	}, tools.ChatRoomSummary)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "AIDrawing",
		Title:       "AI文生图",
		Description: "AI绘图工具，当用户想通过文本生成图像时，可以调用该工具。",
	}, tools.Drawing)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "AIImage2Image",
		Title:       "AI图片修改、图生图",
		Description: "图片修改、编辑、图片合成工具，基于输入的一张或多张图片，结合文本提示词生成新的图片。支持图片混合、风格转换、内容合成等多种创作模式。",
	}, tools.Image2Image)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "ImageRecognition",
		Title:       "图像识别",
		Description: "图像识别工具，当用户想知道图片中的内容或者想提取图片中的内容时，假如你自己不具备多模态能力，可以调用该工具。",
	}, tools.ImageRecognition)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "RequestSong",
		Title:       "点歌",
		Description: "点歌工具，当用户想点播歌曲时，可以调用该工具。",
	}, tools.RequestSong)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "TTS",
		Title:       "文本转语音",
		Description: "文本转语音工具，当用户想让你说话、发语音时，或者想将文本内容转换成语音消息时，可以调用该工具。",
	}, tools.TTS)

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
