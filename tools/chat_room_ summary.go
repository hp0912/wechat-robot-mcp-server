package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sashabaranov/go-openai"

	"wechat-robot-mcp-server/model"
	"wechat-robot-mcp-server/repository"
	"wechat-robot-mcp-server/robot_context"
	"wechat-robot-mcp-server/utils"
)

type ChatRoomSummaryInput struct {
	RecentDuration int `json:"recent_duration" jsonschema:"最近多久的聊天记录，比如总结最近一个小时的聊天记录、总结最近一天的聊天记录。你需要根据用户的需求，转换成秒(示例：最近一小时是3600秒，最近一天是86400秒)。"`
}

func ChatRoomSummary(ctx context.Context, req *mcp.CallToolRequest, params *ChatRoomSummaryInput) (*mcp.CallToolResult, *model.CommonOutput, error) {
	if params.RecentDuration > 24*3600 {
		return utils.CallToolResultError("最多只能总结最近24小时内的聊天记录")
	}

	rc, ok := robot_context.GetRobotContext(ctx)
	if !ok {
		return utils.CallToolResultError("获取机器人上下文失败")
	}

	db, ok := robot_context.GetDB(ctx)
	if !ok {
		return utils.CallToolResultError("获取数据库连接失败")
	}

	globalSettingsRepo := repository.NewGlobalSettingsRepository(ctx, db)
	chatRoomSettingsRepo := repository.NewChatRoomSettingsRepository(ctx, db)
	contactRepo := repository.NewContactRepository(ctx, db)
	messageRepo := repository.NewMessageRepository(ctx, db)

	globalSettings, err := globalSettingsRepo.GetGlobalSettings()
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("获取全局设置失败: %v", err))
	}
	if globalSettings == nil || globalSettings.ChatAIEnabled == nil || !*globalSettings.ChatAIEnabled || globalSettings.ChatAPIKey == "" || globalSettings.ChatBaseURL == "" {
		return utils.CallToolResultError("全局配置群聊总结未开启")
	}

	chatRoomSettings, err := chatRoomSettingsRepo.GetChatRoomSettings(rc.FromWxID)
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("获取群聊设置失败: %v", err))
	}
	if chatRoomSettings == nil || chatRoomSettings.ChatRoomSummaryEnabled == nil || !*chatRoomSettings.ChatRoomSummaryEnabled {
		return utils.CallToolResultError("群聊总结未开启")
	}

	chatRoomName := rc.FromWxID
	chatRoom, err := contactRepo.GetContactByWechatID(rc.FromWxID)
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("获取群聊信息失败: %v", err))
	}
	if chatRoom != nil && chatRoom.Nickname != nil && *chatRoom.Nickname != "" {
		chatRoomName = *chatRoom.Nickname
	}

	startTime := time.Now().Add(-time.Duration(params.RecentDuration) * time.Second)
	endTime := time.Now()
	messages, err := messageRepo.GetMessagesByTimeRange(rc.RobotWxID, rc.FromWxID, startTime.Unix(), endTime.Unix())
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("获取聊天记录失败: %v", err))
	}
	if len(messages) < 100 {
		return utils.CallToolResultError("聊天记录不足100条，不需要总结")
	}

	// 组装对话记录为字符串
	var content []string
	for _, message := range messages {
		// 将时间戳秒格式化为时间YYYY-MM-DD HH:MM:SS 字符串
		timeStr := time.Unix(message.CreatedAt, 0).Format("2006-01-02 15:04:05")
		content = append(content, fmt.Sprintf(`[%s] {"%s": "%s"}--end--`, timeStr, message.Nickname, strings.ReplaceAll(message.Message, "\n", "。。")))
	}
	prompt := `你是一个中文的群聊总结的助手，你可以为一个微信的群聊记录，提取并总结每个时间段大家在重点讨论的话题内容。

每一行代表一个人的发言，每一行的的格式为： {"[time] {nickname}": "{content}"}--end--

请帮我将给出的群聊内容总结成一个今日的群聊报告，包含不多于10个的话题的总结（如果还有更多话题，可以在后面简单补充）。每个话题包含以下内容：
- 话题名(50字以内，带序号1️⃣2️⃣3️⃣，同时附带热度，以🔥数量表示）
- 参与者(不超过5个人，将重复的人名去重)
- 时间段(从几点到几点)
- 过程(50到200字左右）
- 评价(50字以下)
- 分割线： ------------

另外有以下要求：
1. 每个话题结束使用 ------------ 分割
2. 使用中文冒号
3. 无需大标题
4. 开始给出本群讨论风格的整体评价，例如活跃、太水、太黄、太暴力、话题不集中、无聊诸如此类
`
	msg := fmt.Sprintf("群名称: %s\n聊天记录如下:\n%s", chatRoomName, strings.Join(content, "\n"))
	// AI总结
	aiMessages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: prompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: msg,
		},
	}

	// 默认使用AI回复
	aiApiKey := globalSettings.ChatAPIKey
	if *chatRoomSettings.ChatAPIKey != "" {
		aiApiKey = *chatRoomSettings.ChatAPIKey
	}
	aiConfig := openai.DefaultConfig(aiApiKey)
	aiApiBaseURL := strings.TrimRight(globalSettings.ChatBaseURL, "/")
	if chatRoomSettings.ChatBaseURL != nil && *chatRoomSettings.ChatBaseURL != "" {
		aiApiBaseURL = strings.TrimRight(*chatRoomSettings.ChatBaseURL, "/")
	}
	aiConfig.BaseURL = utils.NormalizeAIBaseURL(aiApiBaseURL)
	AIModel := globalSettings.ChatRoomSummaryModel
	if chatRoomSettings.ChatRoomSummaryModel != nil && *chatRoomSettings.ChatRoomSummaryModel != "" {
		AIModel = *chatRoomSettings.ChatRoomSummaryModel
	}
	ai := openai.NewClientWithConfig(aiConfig)
	var resp openai.ChatCompletionResponse
	resp, err = ai.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:               AIModel,
			Messages:            aiMessages,
			Stream:              false,
			MaxCompletionTokens: 2000,
		},
	)
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("AI 总结失败: %v", err))
	}
	// 返回消息为空
	if resp.Choices[0].Message.Content == "" {
		return utils.CallToolResultError("AI 总结失败，返回了空内容")
	}

	replyMsg := fmt.Sprintf("#消息总结\n让我们一起来看看群友们都聊了什么有趣的话题吧~\n\n%s", resp.Choices[0].Message.Content)
	resultContent := []mcp.Content{
		&mcp.TextContent{
			Text: "总结成功",
		},
	}
	output := &model.CommonOutput{
		IsCallToolResult: true,
		ActionType:       model.ActionTypeSendLongTextMessage,
		Text:             replyMsg,
	}

	return &mcp.CallToolResult{
		Content: resultContent,
	}, output, nil
}
