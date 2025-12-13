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
		Name:        "PersonalityAnalysis",
		Title:       "群成员性格分析",
		Description: "根据群聊历史聊天记录，分析群成员的性格特点。",
	}, tools.PersonalityAnalysis)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "AIDrawing",
		Title:       "文生图",
		Description: "AI绘图工具，当用户想通过文本生成图像时，可以调用该工具。",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"prompt": map[string]any{
					"type":        "string",
					"description": "根据用户输入内容，提取出的画图提示词，但是不要对提示词进行总结。",
				},
				"model": map[string]any{
					"type":        "string",
					"description": "画图模型选择（可选）：jimeng-4.0 / jimeng-4.1 / jimeng-4.5，默认 jimeng-4.1。",
					"enum":        []any{"jimeng-4.0", "jimeng-4.1", "jimeng-4.5"},
					"default":     "jimeng-4.1",
				},
				"negative_prompt": map[string]any{
					"type":        "string",
					"description": "用于描述图像中不希望出现的元素或特征的文本，可选。",
				},
				"ratio": map[string]any{
					"type":        "string",
					"description": "图像的宽高比，可选，默认16:9。",
					"default":     "16:9",
				},
				"resolution": map[string]any{
					"type":        "string",
					"description": "图像的分辨率，可选，默认2k。",
					"default":     "2k",
				},
			},
			"required":             []any{"prompt"},
			"additionalProperties": false,
		},
	}, tools.Drawing)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "AIImage2Image",
		Title:       "图片修改、图生图",
		Description: "图片修改、编辑、图片合成工具，基于输入的一张或多张图片，结合文本提示词生成新的图片。支持图片混合、风格转换、内容合成等多种创作模式。输入是文字或图片或两者的组合，输出是图片。",
	}, tools.Image2Image)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "AIVideoGeneration",
		Title:       "视频生成",
		Description: "生成一段视频，支持三种模式：1. 纯文本提示词，不使用任何图片。2. 使用单张图片作为首帧。3. 使用两张图片分别作为首帧和尾帧。",
	}, tools.Image2Video)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "ImageRecognition",
		Title:       "图像识别",
		Description: "图像识别工具，当用户想识别图片中的内容、提取图片中的信息时，可以调用该工具，输入是图片，输出是文字。",
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
	mcp.AddTool(server, &mcp.Tool{
		Name:        "EmojiTool",
		Title:       "表情包工具",
		Description: "当用户想下载表情包、提取表情包时，可以调用该工具。该工具无需参数（调用时 arguments 传 {}）。",
		InputSchema: map[string]any{
			"type":                 "object",
			"properties":           map[string]any{},
			"required":             []any{},
			"additionalProperties": false,
		},
	}, tools.EmojiTool)

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
