package protobuf

import "encoding/xml"

type AppMessage struct {
	XMLName           xml.Name      `xml:"appmsg"`
	AppID             string        `xml:"appid,attr"`
	SDKVer            string        `xml:"sdkver,attr"`
	Title             string        `xml:"title"`
	Des               string        `xml:"des"`
	Action            string        `xml:"action"`
	Type              int           `xml:"type"`
	ShowType          int           `xml:"showtype"`
	Content           string        `xml:"content"`
	URL               string        `xml:"url"`
	DataURL           string        `xml:"dataurl,omitempty"`
	LowURL            string        `xml:"lowurl,omitempty"`
	LowDataURL        string        `xml:"lowdataurl,omitempty"`
	RecordItem        string        `xml:"recorditem"`
	ThumbURL          string        `xml:"thumburl"`
	MessageAction     string        `xml:"messageaction"`
	LanInfo           string        `xml:"laninfo"`
	ExtInfo           string        `xml:"extinfo"`
	SourceUserName    string        `xml:"sourceusername"`
	SourceDisplayName string        `xml:"sourcedisplayname"`
	SongLyric         string        `xml:"songlyric,omitempty"`
	CommentURL        string        `xml:"commenturl"`
	AppAttach         AppAttach     `xml:"appattach"`
	WebViewShared     WebViewShared `xml:"webviewshared"`
	WeAppInfo         WeAppInfo     `xml:"weappinfo"`
	WebSearch         string        `xml:"websearch"`
	SongAlbumURL      string        `xml:"songalbumurl,omitempty"`
}

type AppAttach struct {
	TotalLen    int    `xml:"totallen"`
	AttachID    string `xml:"attachid"`
	EmoticonMD5 string `xml:"emoticonmd5"`
	FileExt     string `xml:"fileext"`
	AesKey      string `xml:"aeskey"`
}

type WebViewShared struct {
	PublisherID    string `xml:"publisherId"`
	PublisherReqID int    `xml:"publisherReqId"`
}

type WeAppInfo struct {
	PagePath       string `xml:"pagepath"`
	UserName       string `xml:"username"`
	AppID          string `xml:"appid"`
	AppServiceType int    `xml:"appservicetype"`
}
