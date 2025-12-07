package protobuf

import "encoding/xml"

type EmojiMsgWrapper struct {
	XMLName xml.Name `xml:"msg"`
	Emoji   EmojiMsg `xml:"emoji"`
}

type EmojiMsg struct {
	FromUserName      string `xml:"fromusername,attr"`
	ToUserName        string `xml:"tousername,attr"`
	Type              int    `xml:"type,attr"`
	IDBuffer          string `xml:"idbuffer,attr"`
	MD5               string `xml:"md5,attr"`
	Len               int    `xml:"len,attr"`
	ProductID         string `xml:"productid,attr"`
	AndroidMD5        string `xml:"androidmd5,attr"`
	AndroidLen        int    `xml:"androidlen,attr"`
	S60v3MD5          string `xml:"s60v3md5,attr"`
	S60v3Len          int    `xml:"s60v3len,attr"`
	S60v5MD5          string `xml:"s60v5md5,attr"`
	S60v5Len          int    `xml:"s60v5len,attr"`
	CDNUrl            string `xml:"cdnurl,attr"`
	DesignerID        string `xml:"designerid,attr"`
	ThumbUrl          string `xml:"thumburl,attr"`
	EncryptUrl        string `xml:"encrypturl,attr"`
	AESKey            string `xml:"aeskey,attr"`
	ExternUrl         string `xml:"externurl,attr"`
	ExternMD5         string `xml:"externmd5,attr"`
	Width             int    `xml:"width,attr"`
	Height            int    `xml:"height,attr"`
	TpUrl             string `xml:"tpurl,attr"`
	TpAuthKey         string `xml:"tpauthkey,attr"`
	AttachedText      string `xml:"attachedtext,attr"`
	AttachedTextColor string `xml:"attachedtextcolor,attr"`
	LensID            string `xml:"lensid,attr"`
	EmojiAttr         string `xml:"emojiattr,attr"`
	LinkID            string `xml:"linkid,attr"`
	Desc              string `xml:"desc,attr"`
}
