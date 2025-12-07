package tools

import (
	"context"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"wechat-robot-mcp-server/model"
	"wechat-robot-mcp-server/protobuf"
	"wechat-robot-mcp-server/repository"
	"wechat-robot-mcp-server/robot_context"
	"wechat-robot-mcp-server/utils"
)

type EmojiInput struct{}

func EmojiTool(ctx context.Context, req *mcp.CallToolRequest, params *EmojiInput) (*mcp.CallToolResult, *model.CommonOutput, error) {
	rc, ok := robot_context.GetRobotContext(ctx)
	if !ok {
		return utils.CallToolResultError("获取机器人上下文失败")
	}

	db, ok := robot_context.GetDB(ctx)
	if !ok {
		return utils.CallToolResultError("获取数据库连接失败")
	}

	messageRepo := repository.NewMessageRepository(ctx, db)
	message, err := messageRepo.GetByID(rc.RefMessageID)
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("查找表情消息失败，你是不是忘了指定需要提取的表情包啊？%v", err))
	}
	if message == nil {
		return utils.CallToolResultError("查找表情消息失败，指定的消息不存在，你是不是忘了指定需要提取的表情包啊？")
	}

	if message.Type == 47 {
		emojiXml := &protobuf.EmojiMsgWrapper{}
		decoder := xml.NewDecoder(strings.NewReader(message.Content))
		err := decoder.Decode(emojiXml)
		if err != nil {
			return utils.CallToolResultError("解析表情消息XML失败")
		}
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: "成功",
					},
				},
			}, &model.CommonOutput{
				IsCallToolResult: true,
				ActionType:       model.ActionTypeEmoji,
				Text:             emojiXml.Emoji.CDNUrl,
			}, nil
	}

	if message.Type == 49 {
		emojiXml := &protobuf.AppMsgWrapper{}
		decoder := xml.NewDecoder(strings.NewReader(message.Content))
		err := decoder.Decode(emojiXml)
		if err != nil {
			return utils.CallToolResultError("解析表情消息XML失败")
		}
		emoji := emojiXml.AppMsg.AppAttach.EmojiInfo
		if emoji == "" {
			return utils.CallToolResultError("指定的应用消息不是表情消息，你是不是忘了指定需要提取的表情包啊？")
		}
		decodedBytes, err := base64.StdEncoding.DecodeString(emoji)
		if err != nil {
			return utils.CallToolResultError(fmt.Sprintf("base64解码失败: %v", err))
		}
		decodedStr := string(decodedBytes)
		re := regexp.MustCompile(`https?://.*?\*`)
		match := re.FindString(decodedStr)
		if match == "" {
			return utils.CallToolResultError("未能从解码后的字符串中提取URL")
		}
		emojiUrl := strings.TrimSuffix(match, "*")
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: "成功",
					},
				},
			}, &model.CommonOutput{
				IsCallToolResult: true,
				ActionType:       model.ActionTypeEmoji,
				Text:             emojiUrl,
			}, nil
	}

	return utils.CallToolResultError("指定的消息不是表情消息，你是不是忘了指定需要提取的表情包啊？")
}
