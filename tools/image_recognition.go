package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sashabaranov/go-openai"

	"wechat-robot-mcp-server/model"
	"wechat-robot-mcp-server/repository"
	"wechat-robot-mcp-server/robot_context"
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

	var aiModel string
	var aiApiKey string
	var aiApiBaseURL string

	globalSettingsRepo := repository.NewGlobalSettingsRepository(ctx, db)
	globalSettings, err := globalSettingsRepo.GetGlobalSettings()
	if err != nil {
		return utils.CallToolResultError("获取全局设置失败")
	}

	if strings.HasSuffix(rc.FromWxID, "@chatroom") {
		chatRoomSettingsRepo := repository.NewChatRoomSettingsRepository(ctx, db)
		chatRoomSettings, err := chatRoomSettingsRepo.GetChatRoomSettings(rc.FromWxID)
		if err != nil {
			return utils.CallToolResultError("获取群聊设置失败")
		}
		if chatRoomSettings != nil {
			if chatRoomSettings.ChatAPIKey != nil {
				aiApiKey = *chatRoomSettings.ChatAPIKey
			}
			if chatRoomSettings.ChatBaseURL != nil {
				aiApiBaseURL = *chatRoomSettings.ChatBaseURL
			}
			if chatRoomSettings.ImageRecognitionModel != nil {
				aiModel = *chatRoomSettings.ImageRecognitionModel
			}
		} else if globalSettings != nil {
			aiApiKey = globalSettings.ChatAPIKey
			aiApiBaseURL = globalSettings.ChatBaseURL
			aiModel = globalSettings.ImageRecognitionModel
		}
	} else {
		friendSettingsRepo := repository.NewFriendSettingsRepo(ctx, db)
		friendSettings, err := friendSettingsRepo.GetFriendSettings(rc.FromWxID)
		if err != nil {
			return utils.CallToolResultError("获取好友设置失败")
		}
		if friendSettings != nil {
			if friendSettings.ChatAPIKey != nil {
				aiApiKey = *friendSettings.ChatAPIKey
			}
			if friendSettings.ChatBaseURL != nil {
				aiApiBaseURL = *friendSettings.ChatBaseURL
			}
			if friendSettings.ImageRecognitionModel != nil {
				aiModel = *friendSettings.ImageRecognitionModel
			}
		} else if globalSettings != nil {
			aiApiKey = globalSettings.ChatAPIKey
			aiApiBaseURL = globalSettings.ChatBaseURL
			aiModel = globalSettings.ImageRecognitionModel
		}
	}

	if aiApiBaseURL == "" || aiApiKey == "" || aiModel == "" {
		return utils.CallToolResultError("AI图片识别未配置，请联系管理员进行配置")
	}

	aiConfig := openai.DefaultConfig(aiApiKey)
	aiConfig.BaseURL = utils.NormalizeAIBaseURL(strings.TrimRight(aiApiBaseURL, "/"))
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
