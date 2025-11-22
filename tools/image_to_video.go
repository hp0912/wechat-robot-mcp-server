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

type Image2VideoInput struct {
	Prompt     string   `json:"prompt" jsonschema:"根据用户输入的文本内容，提取出“生成视频”的提示词，但是不要对提示词进行修改。"`
	FilePaths  []string `json:"file_paths,omitempty" jsonschema:"用于视频的首尾帧的图片地址列表，可选。不提供则表示通过文本生成视频。"`
	Ratio      string   `json:"ratio,omitempty" jsonschema:"生成视频比例，可选。"`
	Resolution string   `json:"resolution,omitempty" jsonschema:"生成视频分辨率，可选。"`
	Duration   int      `json:"duration,omitempty" jsonschema:"生成视频时长，单位秒，可选。"`
}

func Image2Video(ctx context.Context, req *mcp.CallToolRequest, params *Image2VideoInput) (*mcp.CallToolResult, *model.CommonOutput, error) {
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
		// 和图像生成共用 AI 绘图开关
		return utils.CallToolResultError("AI 生成视频未开启")
	}

	aiConfig := settings.GetAIConfig()

	switch aiConfig.ImageModel {
	case model.ImageModelDoubao:
		// Handle 豆包模型
		return utils.CallToolResultError("豆包生成视频暂未实现")
	case model.ImageModelJimeng:
		// Handle 即梦模型
		var jimengConfig pkg.JimengConfig
		if err := json.Unmarshal(aiConfig.ImageAISettings, &jimengConfig); err != nil {
			errmsg := fmt.Sprintf("反序列化即梦绘图配置失败: %v", err)
			log.Print(errmsg)
			return utils.CallToolResultError(errmsg)
		}

		jimengConfig.Prompt = params.Prompt

		var filePaths []*string
		for _, path := range params.FilePaths {
			p := path
			filePaths = append(filePaths, &p)
		}
		jimengConfig.FilePaths = filePaths

		if params.Ratio == "" {
			params.Ratio = "4:3"
		}
		jimengConfig.Ratio = params.Ratio
		// 节约成本，写死 720p
		params.Resolution = "720p"
		jimengConfig.Resolution = params.Resolution
		// 节约成本，只生成 5 秒视频
		params.Duration = 5
		jimengConfig.Duration = params.Duration
		jimengConfig.ResponseFormat = "url"
		imageURLs, err = pkg.JimengVideoGenerations(&jimengConfig)
		if err != nil {
			errmsg := fmt.Sprintf("调用即梦生成视频接口失败: %v", err)
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
