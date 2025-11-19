package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"gorm.io/datatypes"

	"wechat-robot-mcp-server/model"
	"wechat-robot-mcp-server/pkg"
	"wechat-robot-mcp-server/repository"
	"wechat-robot-mcp-server/robot_context"
	"wechat-robot-mcp-server/utils"
)

type DrawingInput struct {
	Prompt         string `json:"prompt" jsonschema:"根据用户输入内容，提取出的画图提示词，但是不要对提示词进行修改。"`
	NegativePrompt string `json:"negative_prompt,omitempty" jsonschema:"用于描述图像中不希望出现的元素或特征的文本，可选。"`
	Ratio          string `json:"ratio,omitempty" jsonschema:"图像的宽高比，可选，默认16:9。"`
	Resolution     string `json:"resolution,omitempty" jsonschema:"图像的分辨率，可选，默认2k。"`
}

func Drawing(ctx context.Context, req *mcp.CallToolRequest, params *DrawingInput) (*mcp.CallToolResult, any, error) {
	rc, ok := robot_context.GetRobotContext(ctx)
	if !ok {
		return utils.CallToolResultError("获取机器人上下文失败")
	}

	db, ok := robot_context.GetDB(ctx)
	if !ok {
		return utils.CallToolResultError("获取数据库连接失败")
	}

	var isDrawingEnabled bool
	var imageModel model.ImageModel = model.ImageModelJimeng
	var imageAISettings datatypes.JSON
	var imageURLs []*string

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
		if chatRoomSettings != nil && chatRoomSettings.ImageAIEnabled != nil {
			isDrawingEnabled = *chatRoomSettings.ImageAIEnabled
			if chatRoomSettings.ImageModel != nil {
				imageModel = *chatRoomSettings.ImageModel
			}
			imageAISettings = chatRoomSettings.ImageAISettings
		} else if globalSettings != nil && globalSettings.ImageAIEnabled != nil {
			isDrawingEnabled = *globalSettings.ImageAIEnabled
			imageModel = globalSettings.ImageModel
			imageAISettings = globalSettings.ImageAISettings
		}
	} else {
		friendSettingsRepo := repository.NewFriendSettingsRepo(ctx, db)
		friendSettings, err := friendSettingsRepo.GetFriendSettings(rc.FromWxID)
		if err != nil {
			return utils.CallToolResultError("获取好友设置失败")
		}
		if friendSettings != nil && friendSettings.ImageAIEnabled != nil {
			isDrawingEnabled = *friendSettings.ImageAIEnabled
			if friendSettings.ImageModel != nil {
				imageModel = *friendSettings.ImageModel
			}
			imageAISettings = friendSettings.ImageAISettings
		} else if globalSettings != nil && globalSettings.ImageAIEnabled != nil {
			isDrawingEnabled = *globalSettings.ImageAIEnabled
			imageModel = globalSettings.ImageModel
			imageAISettings = globalSettings.ImageAISettings
		}
	}

	if !isDrawingEnabled {
		return utils.CallToolResultError("绘图功能未启用")
	}

	switch imageModel {
	case model.ImageModelDoubao:
		// Handle 豆包模型
		var doubaoConfig pkg.DoubaoConfig
		if err := json.Unmarshal(imageAISettings, &doubaoConfig); err != nil {
			errmsg := fmt.Sprintf("反序列化豆包绘图配置失败: %v", err)
			log.Print(errmsg)
			return utils.CallToolResultError(errmsg)
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
		if err := json.Unmarshal(imageAISettings, &jimengConfig); err != nil {
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
		imageURLs, err = pkg.JimengDrawing(&jimengConfig)
		if err != nil {
			errmsg := fmt.Sprintf("调用即梦绘图接口失败: %v", err)
			log.Print(errmsg)
			return utils.CallToolResultError(errmsg)
		}
	case model.ImageModelGLM:
		// Handle 智谱模型
		var glmConfig pkg.GLMConfig
		if err := json.Unmarshal(imageAISettings, &glmConfig); err != nil {
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
		if err := json.Unmarshal(imageAISettings, &hunyuanConfig); err != nil {
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
		return utils.CallToolResultError("获取数据库连接失败")
	case model.ImageModelMidjourney:
		// Handle Midjourney 模型
		return utils.CallToolResultError("获取数据库连接失败")
	case model.ImageModelOpenAI:
		// Handle OpenAI 模型
		return utils.CallToolResultError("获取数据库连接失败")
	default:
		return utils.CallToolResultError("不支持的 AI 图像模型")
	}

	if len(imageURLs) == 0 {
		errmsg := "未生成任何图像"
		log.Print(errmsg)
		return utils.CallToolResultError(errmsg)
	}

	var respData model.BaseResponse
	client := resty.New()
	robotResp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"to_wxid":    rc.FromWxID,
			"image_urls": imageURLs,
		}).
		SetResult(&respData).
		Post(fmt.Sprintf("http://client_%s:%s/api/v1/robot/message/send/image/url", rc.RobotCode, rc.WeChatClientPort))
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("发送图片失败: %v", err))
	}
	if robotResp.StatusCode() != http.StatusOK {
		return utils.CallToolResultError(fmt.Sprintf("发送图片失败，返回状态码不是 200: %d", robotResp.StatusCode()))
	}
	if respData.Code != 200 {
		return utils.CallToolResultError(fmt.Sprintf("发送图片失败，返回状态码不是 200: %s", respData.Message))
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: "绘图结果已发送，你喜欢吗？",
			},
		},
	}, nil, nil
}
