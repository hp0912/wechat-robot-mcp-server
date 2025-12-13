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

type DrawingInput struct {
	Prompt         string `json:"prompt" jsonschema:"根据用户输入内容，提取出的画图提示词，但是不要对提示词进行总结。"`
	Model          string `json:"model,omitempty" jsonschema:"enum=jimeng-4.0,enum=jimeng-4.1,enum=jimeng-4.5,default=jimeng-4.1,画图模型选择"`
	NegativePrompt string `json:"negative_prompt,omitempty" jsonschema:"用于描述图像中不希望出现的元素或特征的文本，可选。"`
	Ratio          string `json:"ratio,omitempty" jsonschema:"图像的宽高比，可选，默认16:9。"`
	Resolution     string `json:"resolution,omitempty" jsonschema:"图像的分辨率，可选，默认2k。"`
}

func Drawing(ctx context.Context, req *mcp.CallToolRequest, params *DrawingInput) (*mcp.CallToolResult, *model.CommonOutput, error) {
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
		var doubaoConfig pkg.DoubaoConfig
		if err := json.Unmarshal(aiConfig.ImageAISettings, &doubaoConfig); err != nil {
			errmsg := fmt.Sprintf("反序列化豆包绘图配置失败: %v", err)
			log.Print(errmsg)
			return utils.CallToolResultError(errmsg)
		}
		if params.Model != "" {
			doubaoConfig.Model = params.Model
		}
		doubaoConfig.Prompt = params.Prompt
		imageURLs, err = pkg.DoubaoDrawing(&doubaoConfig)
		if err != nil {
			errmsg := fmt.Sprintf("调用豆包绘图接口失败: %v", err)
			log.Print(errmsg)
			return utils.CallToolResultError(errmsg)
		}
	case model.ImageModelJimeng:
		// Handle 即梦模型
		var jimengConfig pkg.JimengConfig
		if err := json.Unmarshal(aiConfig.ImageAISettings, &jimengConfig); err != nil {
			errmsg := fmt.Sprintf("反序列化即梦绘图配置失败: %v", err)
			log.Print(errmsg)
			return utils.CallToolResultError(errmsg)
		}

		jimengConfig.Prompt = params.Prompt
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
		imageURLs, err = pkg.JimengImageGenerations(&jimengConfig)
		if err != nil {
			errmsg := fmt.Sprintf("调用即梦绘图接口失败: %v", err)
			log.Print(errmsg)
			return utils.CallToolResultError(errmsg)
		}
	case model.ImageModelGLM:
		// Handle 智谱模型
		var glmConfig pkg.GLMConfig
		if err := json.Unmarshal(aiConfig.ImageAISettings, &glmConfig); err != nil {
			errmsg := fmt.Sprintf("反序列化智谱绘图配置失败: %v", err)
			log.Print(errmsg)
			return utils.CallToolResultError(errmsg)
		}
		glmConfig.Prompt = params.Prompt
		imageURLs, err = pkg.GLMDrawing(&glmConfig)
		if err != nil {
			errmsg := fmt.Sprintf("调用智谱绘图接口失败: %v", err)
			log.Print(errmsg)
			return utils.CallToolResultError(errmsg)
		}
	case model.ImageModelHunyuan:
		// Handle 混元模型
		var hunyuanConfig pkg.HunyuanConfig
		if err := json.Unmarshal(aiConfig.ImageAISettings, &hunyuanConfig); err != nil {
			errmsg := fmt.Sprintf("反序列化混元绘图配置失败: %v", err)
			log.Print(errmsg)
			return utils.CallToolResultError(errmsg)
		}
		hunyuanConfig.Prompt = params.Prompt
		imageURLs, err = pkg.SubmitHunyuanDrawing(&hunyuanConfig)
		if err != nil {
			errmsg := fmt.Sprintf("调用混元绘图接口失败: %v", err)
			log.Print(errmsg)
			return utils.CallToolResultError(errmsg)
		}
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
