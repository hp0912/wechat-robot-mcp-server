package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"wechat-robot-mcp-server/model"
	"wechat-robot-mcp-server/utils"
)

type JoinChatRoomInput struct {
	ChatRoomName string `json:"chat_room_name" jsonschema:"要加入的群聊名称。"`
}

func JoinChatRoom(ctx context.Context, req *mcp.CallToolRequest, params *JoinChatRoomInput) (*mcp.CallToolResult, *model.CommonOutput, error) {
	if params.ChatRoomName == "" {
		return utils.CallToolResultError("群聊名称不能为空")
	}
	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "成功",
				},
			},
		}, &model.CommonOutput{
			IsCallToolResult: true,
			ActionType:       model.ActionTypeJoinChatRoom,
			Text:             params.ChatRoomName,
		}, nil
}
