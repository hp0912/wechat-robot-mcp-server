package tools

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"wechat-robot-mcp-server/config"
	"wechat-robot-mcp-server/model"
	"wechat-robot-mcp-server/protobuf"
	"wechat-robot-mcp-server/utils"
)

type RequestSongInput struct {
	SongTitle string `json:"song_title" jsonschema:"歌曲标题。"`
	Singer    string `json:"singer,omitempty" jsonschema:"歌手。"`
}

type MusicListResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result struct {
		Songs []MusicInfo `json:"songs"`
	} `json:"result"`
}

type MusicInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Al   struct {
		ID     int64  `json:"id"`
		Name   string `json:"name"`
		PicUrl string `json:"picUrl"`
	} `json:"al"`
	Ar []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"ar"`
	Alia []string `json:"alia"`
}

type MusicSearchResponse struct {
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
	Data []MusicSearchData `json:"data"`
}

type MusicSearchData struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

func RequestSong(ctx context.Context, req *mcp.CallToolRequest, params *RequestSongInput) (*mcp.CallToolResult, *model.CommonOutput, error) {
	var result MusicSearchData
	var music MusicInfo
	var list MusicListResponse
	_, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("cookie", fmt.Sprintf("MUSIC_U=%s", config.MUSIC_U)).
		SetQueryParam("keywords", fmt.Sprintf("%s %s", params.SongTitle, params.Singer)).
		SetQueryParam("type", "1").
		SetResult(&list).
		Get("http://netease-cloud-music:3000/search")
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("获取歌曲失败: %v", err))
	}
	if len(list.Result.Songs) == 0 {
		return utils.CallToolResultError(fmt.Sprintf("没有搜索到歌曲 %s", params.SongTitle))
	}

	music = list.Result.Songs[0]

	var resp MusicSearchResponse
	_, err = resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("cookie", fmt.Sprintf("MUSIC_U=%s", config.MUSIC_U)).
		SetQueryParam("id", fmt.Sprintf("%d", music.ID)).
		SetResult(&resp).
		Get("http://netease-cloud-music:3000/song/url")
	if err != nil {
		return utils.CallToolResultError(fmt.Sprintf("获取歌曲失败: %v", err))
	}
	if len(resp.Data) == 0 {
		return utils.CallToolResultError(fmt.Sprintf("没有搜索到歌曲 %s", params.SongTitle))
	}
	result = resp.Data[0]

	musicMsg := protobuf.AppMessage{
		AppID:        "wx8dd6ecd81906fd84",
		SDKVer:       "0",
		Title:        music.Name,
		Action:       "view",
		Type:         76,
		ShowType:     0,
		ThumbURL:     music.Al.PicUrl,
		SongAlbumURL: music.Al.PicUrl,
		URL:          result.URL,
		DataURL:      result.URL,
		LowURL:       result.URL,
		LowDataURL:   result.URL,
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
	if len(music.Alia) > 0 {
		musicMsg.Des = music.Alia[0]
	} else {
		musicMsg.Des = music.Name
	}
	if music.Al.Name != "" {
		musicMsg.Des = fmt.Sprintf("%s - %s", musicMsg.Des, music.Al.Name)
	}
	if len(music.Ar) > 0 {
		artistNames := make([]string, len(music.Ar))
		for i, ar := range music.Ar {
			artistNames[i] = ar.Name
		}
		musicMsg.Des = fmt.Sprintf("%s - %s", musicMsg.Des, strings.Join(artistNames, " "))
	}

	xmlBytes, err := xml.MarshalIndent(musicMsg, "", "  ")
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
