package model

type ActionType int

const (
	ActionTypeSendTextMessage     ActionType = 1 // 发送普通文本消息
	ActionTypeSendLongTextMessage ActionType = 2 // 发送长文本消息
	ActionTypeSendImageMessage    ActionType = 3 // 发送图片消息
	ActionTypeSendVideoMessage    ActionType = 4 // 发送视频消息
	ActionTypeSendAttachMessage   ActionType = 5 // 发送附件消息
	ActionTypeSendVoiceMessage    ActionType = 6 // 发送语音消息
	ActionTypeSendAppMessage      ActionType = 7 // 发送应用消息
	ActionTypeSendEmoticonMessage ActionType = 8 // 发送表情消息
)

type CommonOutput struct {
	IsCallToolResult  bool       `json:"is_call_tool_result,omitempty" jsonschema:"是否为调用工具结果"`
	ActionType        ActionType `json:"action_type" jsonschema:"操作类型: 1-发送普通文本消息, 2-发送长文本消息, 3-发送图片消息, 4-发送视频消息, 5-发送附件消息, 6-发送语音消息, 7-发送应用消息, 8-发送表情消息"`
	Text              string     `json:"text,omitempty" jsonschema:"文本消息内容"`
	AppType           int        `json:"app_type,omitempty" jsonschema:"应用消息类型"`
	AppXML            string     `json:"app_xml,omitempty" jsonschema:"应用消息的XML内容"`
	VoiceEncoding     string     `json:"voice_encoding,omitempty" jsonschema:"语音消息的编码格式"`
	AttachmentURLList []string   `json:"attachment_url_list,omitempty" jsonschema:"附件消息的URL"`
}

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
