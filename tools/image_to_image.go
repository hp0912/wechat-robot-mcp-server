package tools

import (
	"context"
	"encoding/json"
	"fmt"
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
	Model          string   `json:"model,omitempty" jsonschema:"画图模型选择（可选）"`
	Images         []string `json:"images" jsonschema:"用于图片编辑、图片混合、风格转换、内容合成等等的图片链接列表，至少需要一张图像。"`
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

	if params.Model == "" || params.Model == "none" {
		params.Model = "jimeng-5.0"
	}

	switch params.Model {
	case "doubao-seedream-4.5", "doubao-seedream-4.0", "doubao-seedream-3.0-t2i", "doubao-seededit-3.0-i2i":
		// Handle 豆包模型
		return utils.CallToolResultError("豆包图生图暂未实现")
	case "jimeng-4.5", "jimeng-4.6", "jimeng-5.0":
		// Handle 即梦模型
		var config struct {
			JiMeng pkg.JimengConfig `json:"JiMeng"`
		}
		if err := json.Unmarshal(aiConfig.ImageAISettings, &config); err != nil {
			errmsg := fmt.Sprintf("反序列化即梦绘图配置失败: %v", err)
			return utils.CallToolResultError(errmsg)
		}
		if !config.JiMeng.Enabled {
			return utils.CallToolResultError("即梦绘图未开启")
		} else {
			config.JiMeng.Enabled = false
		}
		if params.Model != "" && params.Model != "none" {
			config.JiMeng.Model = params.Model
		} else {
			config.JiMeng.Model = "jimeng-5.0"
		}

		config.JiMeng.Prompt = params.Prompt

		var images []*string
		for _, img := range params.Images {
			images = append(images, &img)
		}
		config.JiMeng.Images = images

		config.JiMeng.NegativePrompt = params.NegativePrompt
		if params.Ratio == "" {
			params.Ratio = "16:9"
		}
		config.JiMeng.Ratio = params.Ratio
		if params.Resolution == "" {
			params.Resolution = "2k"
		}
		config.JiMeng.Resolution = params.Resolution
		config.JiMeng.ResponseFormat = "url"
		imageURLs, err = pkg.JimengImageCompositions(&config.JiMeng)
		if err != nil {
			errmsg := fmt.Sprintf("调用即梦绘图接口失败: %v", err)
			return utils.CallToolResultError(errmsg)
		}
	case "Z-Image", "Z-Image-Turbo", "Qwen-Image-Edit-2511":
		// Handle 造相模型
		var config struct {
			ZImage pkg.ZImageConfig `json:"Z-Image"`
		}
		if err := json.Unmarshal(aiConfig.ImageAISettings, &config); err != nil {
			errmsg := fmt.Sprintf("反序列化造相绘图配置失败: %v", err)
			return utils.CallToolResultError(errmsg)
		}
		if !config.ZImage.Enabled {
			return utils.CallToolResultError("造相绘图未开启")
		} else {
			config.ZImage.Enabled = false
		}
		if params.Model != "" && params.Model != "none" {
			config.ZImage.Model = params.Model
		} else {
			config.ZImage.Model = "Qwen-Image-Edit-2511"
		}
		config.ZImage.Prompt = params.Prompt
		if len(params.Images) > 0 {
			config.ZImage.ImageURL = params.Images
		}
		imageURLs, err = pkg.ZImageDrawing(&config.ZImage)
		if err != nil {
			errmsg := fmt.Sprintf("调用造相绘图接口失败: %v", err)
			return utils.CallToolResultError(errmsg)
		}
	default:
		return utils.CallToolResultError("不支持的 AI 图像模型")
	}

	if len(imageURLs) == 0 {
		errmsg := "未生成任何图像"
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
