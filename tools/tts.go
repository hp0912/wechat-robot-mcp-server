package tools

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
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
	URL           string        `json:"url"`
	RequestHeader RequestHeader `json:"request_header"`
	RequestBody   RequestBody   `json:"request_body"`
}

type RequestHeader struct {
	XApiAppID                        string `json:"X-Api-App-Id"`
	XApiAccessKey                    string `json:"X-Api-Access-Key"`
	XApiResourceID                   string `json:"X-Api-Resource-Id"`
	XApiRequestID                    string `json:"X-Api-Request-Id,omitempty"`
	XControlRequireUsageTokensReturn string `json:"X-Control-Require-Usage-Tokens-Return,omitempty"`
}

type RequestBody struct {
	User      User      `json:"user"`
	Namespace string    `json:"namespace,omitempty"`
	ReqParams ReqParams `json:"req_params"`
}

type User struct {
	UID string `json:"uid,omitempty"`
}

type ReqParams struct {
	Text        string      `json:"text"`
	Model       string      `json:"model"`
	Speaker     string      `json:"speaker"`
	AudioParams AudioParams `json:"audio_params"`
	XAdditions  Additions   `json:"x-additions"`
	Additions   string      `json:"additions,omitempty"`
}

type AudioParams struct {
	Format       string `json:"format,omitempty"`
	SampleRate   int    `json:"sample_rate,omitempty"`
	BitRate      int    `json:"bit_rate,omitempty"`
	Emotion      string `json:"emotion,omitempty"`
	EmotionScale int    `json:"emotion_scale,omitempty"`
	SpeechRate   int    `json:"speech_rate,omitempty"`
	LoudnessRate int    `json:"loudness_rate,omitempty"`
}

type Additions struct {
	SilenceDuration              int      `json:"silence_duration,omitempty"`
	EnableLanguageDetector       bool     `json:"enable_language_detector,omitempty"`
	DisableMarkdownFilter        bool     `json:"disable_markdown_filter,omitempty"`
	DisableEmojiFilter           bool     `json:"disable_emoji_filter,omitempty"`
	MuteCutRemainMs              string   `json:"mute_cut_remain_ms,omitempty"`
	EnableLatexTn                bool     `json:"enable_latex_tn,omitempty"`
	LatexParser                  string   `json:"latex_parser,omitempty"`
	MaxLengthToFilterParenthesis int      `json:"max_length_to_filter_parenthesis,omitempty"`
	ExplicitLanguage             string   `json:"explicit_language,omitempty"`
	ContextLanguage              string   `json:"context_language,omitempty"`
	UnsupportedCharRatioThresh   float64  `json:"unsupported_char_ratio_thresh,omitempty"`
	AigcWatermark                bool     `json:"aigc_watermark,omitempty"`
	ContextTexts                 []string `json:"context_texts,omitempty"`
}

type DoubaoTTSResponse struct {
	Code     int            `json:"code"`
	Message  string         `json:"message"`
	Data     string         `json:"data"`
	Sentence map[string]any `json:"sentence,omitempty"`
	Usage    map[string]any `json:"usage,omitempty"`
}

type TTSInput struct {
	Content      string   `json:"content" jsonschema:"文本转语音的输入文本。"`
	Emotion      string   `json:"emotion,omitempty" jsonschema:"可选，情感类型"`
	ContextTexts []string `json:"context_texts,omitempty" jsonschema:"可选，语音合成的辅助信息，用于模型对话式合成，能更好的体现语音情感"`
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
	if doubaoConfig.URL == "" {
		return utils.CallToolResultError("语音合成地址不能为空")
	}
	if doubaoConfig.RequestHeader.XApiAppID == "" || doubaoConfig.RequestHeader.XApiAccessKey == "" || doubaoConfig.RequestHeader.XApiResourceID == "" {
		return utils.CallToolResultError("请求头参数不能为空")
	}

	doubaoConfig.RequestBody.User.UID = uuid.NewString()
	if doubaoConfig.RequestBody.ReqParams.Speaker == "" {
		doubaoConfig.RequestBody.ReqParams.Speaker = "zh_female_vv_uranus_bigtts"
	}
	doubaoConfig.RequestBody.ReqParams.AudioParams.Format = "mp3"
	doubaoConfig.RequestBody.ReqParams.AudioParams.SampleRate = 24000
	if params.Emotion != "" {
		doubaoConfig.RequestBody.ReqParams.AudioParams.Emotion = params.Emotion
		doubaoConfig.RequestBody.ReqParams.AudioParams.EmotionScale = 5
	}
	if len(params.ContextTexts) > 0 {
		doubaoConfig.RequestBody.ReqParams.XAdditions.ContextTexts = params.ContextTexts
	}
	doubaoConfig.RequestBody.ReqParams.Text = params.Content

	// 准备请求体
	requestBody, err := json.Marshal(doubaoConfig.RequestBody)
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("序列化请求体失败: %v", err))
	}
	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", doubaoConfig.URL, bytes.NewBuffer(requestBody))
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("创建请求失败: %v", err))
	}
	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Api-App-Id", doubaoConfig.RequestHeader.XApiAppID)
	httpReq.Header.Set("X-Api-Access-Key", doubaoConfig.RequestHeader.XApiAccessKey)
	httpReq.Header.Set("X-Api-Resource-Id", doubaoConfig.RequestHeader.XApiResourceID)
	if doubaoConfig.RequestHeader.XApiRequestID != "" {
		httpReq.Header.Set("X-Api-Request-Id", doubaoConfig.RequestHeader.XApiRequestID)
	}
	if doubaoConfig.RequestHeader.XControlRequireUsageTokensReturn != "" {
		httpReq.Header.Set("X-Control-Require-Usage-Tokens-Return", doubaoConfig.RequestHeader.XControlRequireUsageTokensReturn)
	}

	client := &http.Client{Timeout: 300 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("发送请求失败: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return utils.CallToolResultError(fmt.Sprintf("API请求失败，状态码 %d: %s", resp.StatusCode, string(body)))
	}

	audioData := make([]byte, 0)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var ttsResp DoubaoTTSResponse
		if err := json.Unmarshal([]byte(line), &ttsResp); err != nil {
			return utils.CallToolResultError(fmt.Sprintf("解析响应失败: %v, 行内容: %s", err, line))
		}

		if ttsResp.Code == 0 && ttsResp.Data != "" {
			chunkAudio, err := base64.StdEncoding.DecodeString(ttsResp.Data)
			if err != nil {
				return utils.CallToolResultError(fmt.Sprintf("解码音频数据失败: %v", err))
			}
			audioData = append(audioData, chunkAudio...)
			continue
		}

		if ttsResp.Code == 0 && ttsResp.Sentence != nil {
			continue
		}

		// 处理结束标识
		if ttsResp.Code == 20000000 {
			// 合成成功结束
			break
		}

		if ttsResp.Code > 0 {
			return utils.CallToolResultError(fmt.Sprintf("合成失败，错误码: %d, 错误信息: %s", ttsResp.Code, ttsResp.Message))
		}
	}

	if err := scanner.Err(); err != nil {
		return utils.CallToolResultError(fmt.Sprintf("读取响应流失败: %v", err))
	}

	if len(audioData) == 0 {
		return utils.CallToolResultError("未接收到音频数据")
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
			Text:             base64.StdEncoding.EncodeToString(audioData),
			VoiceEncoding:    doubaoConfig.RequestBody.ReqParams.AudioParams.Format,
		}, nil
}
