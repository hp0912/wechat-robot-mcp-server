package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sashabaranov/go-openai"

	"wechat-robot-mcp-server/interface/settings"
	"wechat-robot-mcp-server/model"
	"wechat-robot-mcp-server/robot_context"
	"wechat-robot-mcp-server/service"
	"wechat-robot-mcp-server/utils"
)

type ImageRecognitionInput struct {
	Prompt   string `json:"prompt" jsonschema:"图像识别提示词，用户想对图片做什么处理。"`
	ImageURL string `json:"image_url" jsonschema:"图片的URL地址。"`
}

func ImageRecognition(ctx context.Context, req *mcp.CallToolRequest, params *ImageRecognitionInput) (*mcp.CallToolResult, *model.CommonOutput, error) {
	rc, ok := robot_context.GetRobotContext(ctx)
	if !ok {
		return utils.CallToolResultError("获取机器人上下文失败")
	}

	db, ok := robot_context.GetDB(ctx)
	if !ok {
		return utils.CallToolResultError("获取数据库连接失败")
	}

	var settings settings.Settings
	var err error

	if strings.HasSuffix(rc.FromWxID, "@chatroom") {
		settings = service.NewChatRoomSettingsService(ctx, db)
	} else {
		settings = service.NewFriendSettingsService(ctx, db)
	}
	err = settings.InitByMessage(&model.Message{
		FromWxID: rc.FromWxID,
	})
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("初始化 AI 设置失败: %v", err))
	}

	aiConf := settings.GetAIConfig()
	aiApiKey := aiConf.APIKey
	aiApiBaseURL := aiConf.BaseURL
	aiModel := aiConf.ImageRecognitionModel

	if aiApiBaseURL == "" || aiApiKey == "" || aiModel == "" {
		return utils.CallToolResultError("AI图片识别未配置，请联系管理员进行配置")
	}

	aiConfig := openai.DefaultConfig(aiApiKey)
	aiConfig.BaseURL = aiApiBaseURL
	ai := openai.NewClientWithConfig(aiConfig)
	var resp openai.ChatCompletionResponse
	resp, err = ai.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: aiModel,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleUser, MultiContent: []openai.ChatMessagePart{
					{
						Type: openai.ChatMessagePartTypeImageURL,
						ImageURL: &openai.ChatMessageImageURL{
							URL: params.ImageURL,
						},
					},
					{
						Type: openai.ChatMessagePartTypeText,
						Text: params.Prompt,
					},
				}},
			},
			Stream: false,
		},
	)
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("图片识别失败: %v", err))
	}
	// 返回消息为空
	if resp.Choices[0].Message.Content == "" {
		return utils.CallToolResultError("图片识别失败，返回了空内容")
	}

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: resp.Choices[0].Message.Content,
				},
			},
		}, &model.CommonOutput{
			IsCallToolResult: true,
			ActionType:       model.ActionTypeSendTextMessage,
			Text:             resp.Choices[0].Message.Content,
		}, nil
}
