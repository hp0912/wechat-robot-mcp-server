package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"wechat-robot-mcp-server/interface/settings"
	"wechat-robot-mcp-server/model"
	"wechat-robot-mcp-server/pkg"
	"wechat-robot-mcp-server/robot_context"
	"wechat-robot-mcp-server/service"
	"wechat-robot-mcp-server/utils"
)

type Image2ImageInput struct {
	Prompt         string   `json:"prompt" jsonschema:"根据用户输入的文本内容，提取出图片混合、风格转换、内容合成等等的提示词，但是不要对提示词进行修改。"`
	Images         []string `json:"images" jsonschema:"用于t图片编辑、图片混合、风格转换、内容合成等等的图片链接列表，至少需要一张图像。"`
	NegativePrompt string   `json:"negative_prompt,omitempty" jsonschema:"用于描述图像中不希望出现的元素或特征的文本，可选。"`
	Ratio          string   `json:"ratio,omitempty" jsonschema:"图像的宽高比，可选，默认16:9。"`
	Resolution     string   `json:"resolution,omitempty" jsonschema:"图像的分辨率，可选，默认2k。"`
}

func Image2Image(ctx context.Context, req *mcp.CallToolRequest, params *Image2ImageInput) (*mcp.CallToolResult, *model.CommonOutput, error) {
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
	var imageURLs []*string

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

	if !settings.IsAIDrawingEnabled() {
		return utils.CallToolResultError("AI 绘图未开启")
	}

	aiConfig := settings.GetAIConfig()

	switch aiConfig.ImageModel {
	case model.ImageModelDoubao:
		// Handle 豆包模型
		return utils.CallToolResultError("豆包图生图暂未实现")
	case model.ImageModelJimeng:
		// Handle 即梦模型
		var jimengConfig pkg.JimengConfig
		if err := json.Unmarshal(aiConfig.ImageAISettings, &jimengConfig); err != nil {
			errmsg := fmt.Sprintf("反序列化即梦绘图配置失败: %v", err)
			log.Print(errmsg)
			return utils.CallToolResultError(errmsg)
		}

		jimengConfig.Prompt = params.Prompt
		jimengConfig.Images = params.Images
		jimengConfig.NegativePrompt = params.NegativePrompt
		if params.Ratio == "" {
			params.Ratio = "16:9"
		}
		jimengConfig.Ratio = params.Ratio
		if params.Resolution == "" {
			params.Resolution = "2k"
		}
		jimengConfig.Resolution = params.Resolution
		jimengConfig.ResponseFormat = "url"
		imageURLs, err = pkg.JimengImageCompositions(&jimengConfig)
		if err != nil {
			errmsg := fmt.Sprintf("调用即梦绘图接口失败: %v", err)
			log.Print(errmsg)
			return utils.CallToolResultError(errmsg)
		}
	case model.ImageModelGLM:
		// Handle 智谱模型
	case model.ImageModelHunyuan:
		// Handle 混元模型
	case model.ImageModelStableDiffusion:
		// Handle Stable Diffusion 模型
	case model.ImageModelMidjourney:
		// Handle Midjourney 模型
	case model.ImageModelOpenAI:
		// Handle OpenAI 模型
	default:
		return utils.CallToolResultError("不支持的 AI 图像模型")
	}

	if len(imageURLs) == 0 {
		errmsg := "未生成任何图像"
		log.Print(errmsg)
		return utils.CallToolResultError(errmsg)
	}

	var attachmentURLList []string
	for _, url := range imageURLs {
		if url != nil {
			attachmentURLList = append(attachmentURLList, *url)
		}
	}

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "绘图成功",
				},
			},
		}, &model.CommonOutput{
			IsCallToolResult:  true,
			ActionType:        model.ActionTypeSendImageMessage,
			AttachmentURLList: attachmentURLList,
		}, nil
}
