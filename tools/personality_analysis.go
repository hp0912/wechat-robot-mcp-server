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
	"wechat-robot-mcp-server/service"
	"wechat-robot-mcp-server/utils"
)

type PersonalityAnalysisInput struct {
	RecentDuration int    `json:"recent_duration" jsonschema:"最近多久的聊天记录，比如最近一个小时的聊天记录、最近一天的聊天记录。你需要根据用户的需求，转换成秒(示例：最近一小时是3600秒，最近一天是86400秒)。"`
	Nickname       string `json:"nickname,omitempty" jsonschema:"需要分析的群成员昵称，可选。"`
}

func PersonalityAnalysis(ctx context.Context, req *mcp.CallToolRequest, params *PersonalityAnalysisInput) (*mcp.CallToolResult, *model.CommonOutput, error) {
	if params.RecentDuration > 24*3600 {
		return utils.CallToolResultError("最多只支持最近24小时内的聊天记录")
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
	chatRoomMemberRepo := repository.NewChatRoomMemberRepo(ctx, db)
	contactRepo := repository.NewContactRepository(ctx, db)
	messageRepo := repository.NewMessageRepository(ctx, db)

	globalSettings, err := globalSettingsRepo.GetGlobalSettings()
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("获取全局设置失败: %v", err))
	}
	if globalSettings == nil || globalSettings.ChatAIEnabled == nil || !*globalSettings.ChatAIEnabled || globalSettings.ChatAPIKey == "" || globalSettings.ChatBaseURL == "" {
		return utils.CallToolResultError("全局配置群聊总结未开启，不支持群成员性格分析")
	}

	chatRoomSettings, err := chatRoomSettingsRepo.GetChatRoomSettings(rc.FromWxID)
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("获取群聊设置失败: %v", err))
	}
	if chatRoomSettings == nil || chatRoomSettings.ChatRoomSummaryEnabled == nil || !*chatRoomSettings.ChatRoomSummaryEnabled {
		return utils.CallToolResultError("群聊总结未开启，不支持群成员性格分析")
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
	if len(messages) < 50 {
		return utils.CallToolResultError("聊天记录有点少啊，分析可能会不准确，多聊会儿吧")
	}

	// 组装对话记录为字符串
	var content []string
	for _, message := range messages {
		// 将时间戳秒格式化为时间YYYY-MM-DD HH:MM:SS 字符串
		timeStr := time.Unix(message.CreatedAt, 0).Format("2006-01-02 15:04:05")
		content = append(content, fmt.Sprintf(`[%s] {"%s": "%s"}--end--`, timeStr, message.Nickname, strings.ReplaceAll(message.Message, "\n", "。。")))
	}
	prompt := `
你现在处于一个在线群聊中，群聊中有多个人在发言。你的任务是根据指定用户的发言内容，分析并推断该用户的 MBTI 性格类型。

## 核心目标 (Core Objective)
通过深入分析提供的聊天记录内容，运用 MBTI (Myers-Briggs Type Indicator) 理论框架，准确推断指定用户的性格类型，并提供详细的维度分析和依据。

## 角色与背景 (Role & Context)
你是一位拥有心理学背景的资深人格分析专家，精通 MBTI 理论及其认知功能（Ne, Ni, Se, Si, Te, Ti, Fe, Fi）。你的任务是从日常对话的碎片化信息中，敏锐地捕捉用户的行为模式、思维逻辑、决策方式和能量来源，从而构建出完整的人物性格画像。

## 关键指令与步骤 (Key Instructions & Steps)
1.  **数据预处理**：仔细阅读提供的聊天记录，区分“指定用户”与其他参与者的发言，排除无关噪音。
2.  **维度拆解**：针对 MBTI 的四个维度进行逐一分析：
    *   **E/I (外向/内向)**：分析用户是从外部互动中获取能量，还是通过独处反思恢复精力。
    *   **S/N (实感/直觉)**：观察用户关注的是具体细节、现实当下，还是抽象概念、未来可能性。
    *   **T/F (理智/情感)**：判断用户做决策时依据的是逻辑规则、客观事实，还是个人价值观、他人感受。
    *   **J/P (判断/感知)**：分析用户的生活方式是喜欢计划、结构化，还是偏好灵活、随性。
3.  **认知功能验证**：尝试识别用户使用的主要认知功能（如主导功能和辅助功能），以验证四个维度的判断是否自洽。
4.  **证据引用**：在得出每个维度的结论时，必须直接引用聊天记录中的具体语句作为证据。
5.  **综合结论**：确定最终的四字母 MBTI 类型，并描述该类型的典型特征与用户的匹配度。

## 输入信息 (Input Data / Information)
*   **待分析数据**：用户提供的具体聊天记录文本。
*   **指定对象**：聊天记录中需要被分析的具体用户名或昵称（由用户在后续输入中指定）。

## 输出要求 (Output Requirements)
请按以下 Markdown 格式输出分析报告：

*   **最终判定**：[MBTI 类型代码，例如：INTJ] - [类型名称，例如：建筑师型]
*   **详细维度分析**：
    *   **能量倾向 (E/I)**：结论 + 聊天记录证据分析。
    *   **信息获取 (S/N)**：结论 + 聊天记录证据分析。
    *   **决策方式 (T/F)**：结论 + 聊天记录证据分析。
    *   **生活态度 (J/P)**：结论 + 聊天记录证据分析。
*   **性格画像总结**：用一段话总结该用户的核心性格特征、潜在优势及可能的盲点。
*   **置信度**：对本次分析准确性的自我评估（高/中/低），并说明理由（如样本量是否充足）。

## 约束与偏好 (Constraints & Preferences)
*   **客观中立**：分析必须基于提供的文本证据，避免无根据的猜测或刻板印象。
*   **避免巴纳姆效应**：不要使用模糊、笼统适用于所有人的描述，要具体针对该用户。
*   **语言风格**：保持专业、理性且具有洞察力。
*   **隐私保护**：分析中不要重复提及敏感的个人身份信息（如真实电话号码、地址等）。`

	msg := "### 每一行代表一个人的发言，每一行的的格式为： {\"[time] {nickname}\": \"{content}\"}--end--\n"

	if params.Nickname == "" {
		senders, err := chatRoomMemberRepo.GetChatRoomMemberByWeChatIDs(rc.FromWxID, []string{rc.SenderWxID})
		if err != nil {
			return utils.CallToolResultError(fmt.Sprintf("获取群成员信息失败: %v", err))
		}
		if len(senders) == 0 {
			return utils.CallToolResultError("获取群成员信息失败，指定的成员不存在")
		}
		if senders[0].Remark != "" {
			params.Nickname = senders[0].Remark
		} else {
			params.Nickname = senders[0].Nickname
		}
	}

	msg += fmt.Sprintf("### 着重分析昵称为 '%s' 的群成员的性格类型。\n", params.Nickname)
	msg += "### 简单分析下所有参与聊天人员的性格类型。\n\n"
	msg += fmt.Sprintf("群名称: %s\n聊天记录如下:\n%s", chatRoomName, strings.Join(content, "\n"))
	// AI 分析
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

	settings := service.NewChatRoomSettingsService(context.Background(), db)
	err = settings.InitByMessage(&model.Message{
		FromWxID: rc.FromWxID,
	})
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("初始化群聊设置失败: %v", err))
	}

	aiConf := settings.GetAIConfig()

	aiConfig := openai.DefaultConfig(aiConf.APIKey)
	aiConfig.BaseURL = aiConf.BaseURL
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
		return utils.CallToolResultError(fmt.Sprintf("AI 分析失败: %v", err))
	}
	// 返回消息为空
	if resp.Choices[0].Message.Content == "" {
		return utils.CallToolResultError("AI 分析失败，返回了空内容")
	}

	resultContent := []mcp.Content{
		&mcp.TextContent{
			Text: "分析成功",
		},
	}
	output := &model.CommonOutput{
		IsCallToolResult: true,
		ActionType:       model.ActionTypeSendLongTextMessage,
		Text:             resp.Choices[0].Message.Content,
	}

	return &mcp.CallToolResult{
		Content: resultContent,
	}, output, nil
}
