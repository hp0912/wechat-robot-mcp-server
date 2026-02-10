package tools

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"wechat-robot-mcp-server/model"
	"wechat-robot-mcp-server/protobuf"
	"wechat-robot-mcp-server/utils"
)

type RequestSongInput struct {
	SongTitle string `json:"song_title" jsonschema:"歌曲标题。"`
	Singer    string `json:"singer,omitempty" jsonschema:"歌手。"`
}

type MusicListResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Artist string `json:"artist"`
	} `json:"data"`
}

type MusicSearchResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data MusicSearchData `json:"data"`
}

type MusicSearchData struct {
	ID     string  `json:"id"`
	Name   *string `json:"name"`
	Album  *string `json:"album"`
	Artist string  `json:"artist"`
	Pic    *string `json:"pic"`
	URL    string  `json:"url"`
	Lrc    *string `json:"lrc"`
}

func RequestSong(ctx context.Context, req *mcp.CallToolRequest, params *RequestSongInput) (*mcp.CallToolResult, *model.CommonOutput, error) {
	var result MusicSearchData
	var musicId string
	var list MusicListResponse

	_, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("name", params.SongTitle).
		SetQueryParam("type", "qq").
		SetQueryParam("page", "1").
		SetQueryParam("limit", "20").
		SetResult(&list).
		Get("https://yunzhiapi.cn/API/hqyyid.php")
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("获取歌曲失败: %v", err))
	}
	if len(list.Data) == 0 {
		return utils.CallToolResultError(fmt.Sprintf("没有搜索到歌曲 %s", params.SongTitle))
	}

	if params.Singer != "" {
		for _, item := range list.Data {
			if strings.Contains(item.Artist, params.Singer) {
				musicId = item.ID
				break
			}
		}
		if musicId == "" {
			musicId = list.Data[0].ID
		}
	} else {
		musicId = list.Data[0].ID
	}

	var resp MusicSearchResponse
	_, err = resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("id", musicId).
		SetQueryParam("type", "qq").
		SetResult(&resp).
		Get("https://yunzhiapi.cn/API/yyjhss.php")
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("获取歌曲失败: %v", err))
	}
	result = resp.Data
	if result.Name == nil || result.Album == nil {
		return utils.CallToolResultError(fmt.Sprintf("没有搜索到歌曲 %s", params.SongTitle))
	}

	music := protobuf.AppMessage{
		AppID:      "wx5aa333606550dfd5",
		SDKVer:     "0",
		Title:      *result.Name,
		Des:        fmt.Sprintf("%s -%s", *result.Album, result.Artist),
		Action:     "view",
		Type:       3,
		ShowType:   0,
		URL:        result.URL,
		DataURL:    result.URL,
		LowURL:     result.URL,
		LowDataURL: result.URL,
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
	if result.Pic != nil {
		music.ThumbURL = *result.Pic
		music.SongAlbumURL = *result.Pic
	}

	xmlBytes, err := xml.MarshalIndent(music, "", "  ")
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("序列化歌曲失败: %v", err))
	}

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "点播成功",
				},
			},
		}, &model.CommonOutput{
			IsCallToolResult: true,
			ActionType:       model.ActionTypeSendAppMessage,
			AppType:          3,
			AppXML:           string(xmlBytes),
		}, nil
}
