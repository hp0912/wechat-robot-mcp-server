package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"wechat-robot-mcp-server/interface/settings"
	"wechat-robot-mcp-server/model"
	"wechat-robot-mcp-server/robot_context"
	"wechat-robot-mcp-server/service"
	"wechat-robot-mcp-server/utils"
)

type DoubaoTTSConfig struct {
	BaseURL     string `json:"base_url"`
	AccessToken string `json:"access_token"`
	DoubaoTTSRequest
}

type DoubaoTTSRequest struct {
	App     AppConfig     `json:"app"`
	User    UserConfig    `json:"user"`
	Audio   AudioConfig   `json:"audio"`
	Request RequestConfig `json:"request"`
}

type AppConfig struct {
	AppID   string `json:"appid"`
	Token   string `json:"token"`
	Cluster string `json:"cluster"`
}

type UserConfig struct {
	UID string `json:"uid"`
}

type AudioConfig struct {
	VoiceType       string  `json:"voice_type"`
	Encoding        string  `json:"encoding"`
	CompressionRate int     `json:"compression_rate"`
	Rate            int     `json:"rate"`
	SpeedRatio      float64 `json:"speed_ratio"`
	VolumeRatio     float64 `json:"volume_ratio"`
	PitchRatio      float64 `json:"pitch_ratio"`
	Emotion         string  `json:"emotion"`
	Language        string  `json:"language"`
}

type RequestConfig struct {
	ReqID           string `json:"reqid"`
	Text            string `json:"text"`
	TextType        string `json:"text_type"`
	Operation       string `json:"operation"`
	SilenceDuration string `json:"silence_duration"`
	WithFrontend    string `json:"with_frontend"`
	FrontendType    string `json:"frontend_type"`
	PureEnglishOpt  string `json:"pure_english_opt"`
}

type DoubaoTTSResponse struct {
	ReqID     string   `json:"reqid"`
	Code      int      `json:"code"`
	Operation string   `json:"operation"`
	Message   string   `json:"message"`
	Sequence  int      `json:"sequence"`
	Data      string   `json:"data"`
	Addition  Addition `json:"addition"`
}

type Addition struct {
	Description string `json:"description"`
	Duration    string `json:"duration"`
	Frontend    string `json:"frontend"`
}

type TTSInput struct {
	Content string `json:"content" jsonschema:"文本转语音的输入文本。"`
}

func TTS(ctx context.Context, req *mcp.CallToolRequest, params *TTSInput) (*mcp.CallToolResult, *model.CommonOutput, error) {
	rc, ok := robot_context.GetRobotContext(ctx)
	if !ok {
		return utils.CallToolResultError("获取机器人上下文失败")
	}

	db, ok := robot_context.GetDB(ctx)
	if !ok {
		return utils.CallToolResultError("获取数据库连接失败")
	}

	if params.Content == "" {
		return utils.CallToolResultError("文本转语音的输入文本不能为空")
	}
	if utf8.RuneCountInString(params.Content) > 260 {
		return utils.CallToolResultError("你要说的也太多了，要不你还是说点别的吧。")
	}

	var settings settings.Settings

	if strings.HasSuffix(rc.FromWxID, "@chatroom") {
		settings = service.NewChatRoomSettingsService(ctx, db)
	} else {
		settings = service.NewFriendSettingsService(ctx, db)
	}
	err := settings.InitByMessage(&model.Message{
		FromWxID: rc.FromWxID,
	})
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("初始化 AI 设置失败: %v", err))
	}

	if !settings.IsTTSEnabled() {
		return utils.CallToolResultError("文本转语音未开启")
	}

	aiConfig := settings.GetAIConfig()
	var doubaoConfig DoubaoTTSConfig
	if err := json.Unmarshal(aiConfig.TTSSettings, &doubaoConfig); err != nil {
		return utils.CallToolResultError("反序列化豆包文本转语音配置失败: " + err.Error())
	}
	doubaoConfig.Request.Text = params.Content

	if doubaoConfig.App.AppID == "" {
		return utils.CallToolResultError("应用ID不能为空")
	}
	if doubaoConfig.AccessToken == "" {
		return utils.CallToolResultError("未找到语音合成密钥")
	}

	doubaoConfig.App.Token = uuid.NewString()
	doubaoConfig.App.Cluster = "volcano_tts"
	doubaoConfig.User.UID = uuid.NewString()
	doubaoConfig.Audio.Encoding = "mp3"
	doubaoConfig.Request.ReqID = uuid.NewString()
	doubaoConfig.Request.Operation = "query"
	doubaoConfig.Request.TextType = "plain"

	// 准备请求体
	requestBody, err := json.Marshal(doubaoConfig.DoubaoTTSRequest)
	if err != nil {
		return utils.CallToolResultError(err.Error())
	}
	// 创建HTTP请求
	ttsReq, err := http.NewRequest("POST", doubaoConfig.BaseURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("创建请求失败: %v", err))
	}
	// 设置请求头
	ttsReq.Header.Set("Content-Type", "application/json")
	ttsReq.Header.Set("Authorization", fmt.Sprintf("Bearer; %s", doubaoConfig.AccessToken))

	// 发送请求
	client := &http.Client{Timeout: 300 * time.Second}
	resp, err := client.Do(ttsReq)
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("发送请求失败: %v", err))
	}
	defer resp.Body.Close()
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("读取响应失败: %v", err))
	}
	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return utils.CallToolResultError(fmt.Sprintf("API请求失败，状态码 %d: %s", resp.StatusCode, string(body)))
	}
	// 解析响应
	var ttsResp DoubaoTTSResponse
	if err := json.Unmarshal(body, &ttsResp); err != nil {
		return utils.CallToolResultError(fmt.Sprintf("解析响应失败: %v", err))
	}
	if ttsResp.Message != "Success" {
		return utils.CallToolResultError(fmt.Sprintf("合成失败: %s", ttsResp.Message))
	}

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "成功",
				},
			},
		}, &model.CommonOutput{
			IsCallToolResult: true,
			ActionType:       model.ActionTypeSendVoiceMessage,
			Text:             ttsResp.Data,
			VoiceEncoding:    doubaoConfig.Audio.Encoding,
		}, nil
}
