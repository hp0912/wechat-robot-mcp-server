package tools

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"wechat-robot-mcp-server/model"
	"wechat-robot-mcp-server/protobuf"
	"wechat-robot-mcp-server/robot_context"
	"wechat-robot-mcp-server/utils"
)

type RequestSongInput struct {
	SongTitle string `json:"song_title" jsonschema:"歌曲标题。"`
}

type MusicSearchResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data MusicSearchData `json:"data"`
}

type MusicSearchData struct {
	Title    *string `json:"title"`
	Singer   string  `json:"singer"`
	ID       string  `json:"id"`
	Cover    *string `json:"cover"`
	Link     string  `json:"link"`
	MusicURL string  `json:"music_url"`
	Lrc      *string `json:"lrc"`
}

func RequestSong(ctx context.Context, req *mcp.CallToolRequest, params *RequestSongInput) (*mcp.CallToolResult, any, error) {
	rc, ok := robot_context.GetRobotContext(ctx)
	if !ok {
		return utils.CallToolResultError("获取机器人上下文失败")
	}

	var resp MusicSearchResponse
	_, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("msg", params.SongTitle).
		SetQueryParam("type", "json").
		SetQueryParam("n", "1").
		SetQueryParam("br", "7").
		SetResult(&resp).
		Get("https://api.cenguigui.cn/api/music/netease/WyY_Dg.php")
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("获取歌曲失败: %v", err))
	}
	result := resp.Data
	if result.Title == nil {
		return utils.CallToolResultError(fmt.Sprintf("没有搜索到歌曲 %s", params.SongTitle))
	}

	music := protobuf.AppMessage{
		AppID:      "wx8dd6ecd81906fd84",
		SDKVer:     "0",
		Title:      *result.Title,
		Des:        result.Singer,
		Action:     "view",
		Type:       3,
		ShowType:   0,
		URL:        result.Link,
		DataURL:    result.MusicURL,
		LowURL:     result.Link,
		LowDataURL: result.MusicURL,
		AppAttach: protobuf.AppAttach{
			TotalLen: 0,
		},
		WebViewShared: protobuf.WebViewShared{
			PublisherReqID: 0,
		},
		WeAppInfo: protobuf.WeAppInfo{
			AppServiceType: 0,
		},
	}
	if result.Lrc != nil {
		music.SongLyric = *result.Lrc
	}
	if result.Cover != nil {
		music.ThumbURL = *result.Cover
		music.SongAlbumURL = *result.Cover
	}

	xmlBytes, err := xml.MarshalIndent(music, "", "  ")
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("序列化歌曲失败: %v", err))
	}

	var respData model.BaseResponse
	client := resty.New()
	robotResp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"to_wxid": rc.FromWxID,
			"type":    3,
			"xml":     string(xmlBytes),
		}).
		SetResult(&respData).
		Post(fmt.Sprintf("http://client_%s:%s/api/v1/robot/message/send/app", rc.RobotCode, rc.WeChatClientPort))
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("发送歌曲失败: %v", err))
	}
	if robotResp.StatusCode() != http.StatusOK {
		return utils.CallToolResultError(fmt.Sprintf("发送歌曲失败，返回状态码不是 200: %d", robotResp.StatusCode()))
	}
	if respData.Code != 200 {
		return utils.CallToolResultError(fmt.Sprintf("发送歌曲失败，返回状态码不是 200: %s", respData.Message))
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: "歌曲：" + *result.Title + " - " + result.Singer + " 点播成功",
			},
		},
	}, nil, nil
}
