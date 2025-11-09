package main

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type WelcomeType string

const (
	WelcomeTypeText  WelcomeType = "text"  // 文本
	WelcomeTypeEmoji WelcomeType = "emoji" // 表情
	WelcomeTypeImage WelcomeType = "image" // 图片
	WelcomeTypeURL   WelcomeType = "url"   // 链接
)

type PatType string

const (
	PatTypeText  PatType = "text"  // 文本
	PatTypeVoice PatType = "voice" // 语音
)

type NewsType string

const (
	NewsTypeNone  NewsType = ""
	NewsTypeText  NewsType = "text"  // 文本
	NewsTypeImage NewsType = "image" // 图片
)

type ImageModel string

const (
	ImageModelDoubao          ImageModel = "doubao"           // 豆包模型
	ImageModelJimeng          ImageModel = "jimeng"           // 即梦模型
	ImageModelGLM             ImageModel = "glm"              // 智谱模型
	ImageModelHunyuan         ImageModel = "hunyuan"          // 腾讯混元模型
	ImageModelStableDiffusion ImageModel = "stable-diffusion" // Stable Diffusion 模型
	ImageModelMidjourney      ImageModel = "midjourney"       // Midjourney 模型
	ImageModelOpenAI          ImageModel = "openai"           // OpenAI 模型
)

type GlobalSettings struct {
	ID                        uint64         `gorm:"column:id;primaryKey;autoIncrement;comment:公共配置表主键ID" json:"id"`
	ChatAIEnabled             *bool          `gorm:"column:chat_ai_enabled;default:false;comment:是否启用AI聊天功能" json:"chat_ai_enabled"`
	ChatAITrigger             *string        `gorm:"column:chat_ai_trigger;type:varchar(20);default:'';comment:触发聊天AI的关键词" json:"chat_ai_trigger"`
	ChatBaseURL               string         `gorm:"column:chat_base_url;type:varchar(255);default:'';comment:聊天AI的基础URL地址" json:"chat_base_url"`
	ChatAPIKey                string         `gorm:"column:chat_api_key;type:varchar(255);default:'';comment:聊天AI的API密钥" json:"chat_api_key"`
	WorkflowModel             string         `gorm:"column:workflow_model;type:varchar(100);default:'';comment:聊天AI使用的模型名称" json:"workflow_model"`
	ChatModel                 string         `gorm:"column:chat_model;type:varchar(100);default:'';comment:聊天AI使用的模型名称" json:"chat_model"`
	ImageRecognitionModel     string         `gorm:"column:image_recognition_model;type:varchar(100);default:'';comment:图像识别AI使用的模型名称" json:"image_recognition_model"`
	ChatPrompt                string         `gorm:"column:chat_prompt;type:text;comment:聊天AI系统提示词" json:"chat_prompt"`
	MaxCompletionTokens       *int           `gorm:"column:max_completion_tokens;default:0;comment:最大回复" json:"max_completion_tokens"`
	ImageAIEnabled            *bool          `gorm:"column:image_ai_enabled;default:false;comment:是否启用AI绘图功能" json:"image_ai_enabled"`
	ImageModel                ImageModel     `gorm:"column:image_model;type:varchar(255);default:'';comment:绘图AI模型" json:"image_model"`
	ImageAISettings           datatypes.JSON `gorm:"column:image_ai_settings;type:json;comment:绘图AI配置项" json:"image_ai_settings"`
	TTSEnabled                *bool          `gorm:"column:tts_enabled;default:false;comment:是否启用AI文本转语音功能" json:"tts_enabled"`
	TTSSettings               datatypes.JSON `gorm:"column:tts_settings;type:json;comment:文本转语音配置项" json:"tts_settings"`
	LTTSSettings              datatypes.JSON `gorm:"column:ltts_settings;type:json;comment:长文本转语音配置项" json:"ltts_settings"`
	PatEnabled                *bool          `gorm:"column:pat_enabled;default:false;comment:是否启用拍一拍功能" json:"pat_enabled"`
	PatType                   PatType        `gorm:"column:pat_type;type:enum('text','voice');default:'text';comment:拍一拍方式：text-文本，voice-语音" json:"pat_type"`
	PatText                   string         `gorm:"column:pat_text;type:varchar(255);default:'';comment:拍一拍的文本" json:"pat_text"`
	PatVoiceTimbre            string         `gorm:"column:pat_voice_timbre;type:varchar(255);default:'';comment:拍一拍的音色" json:"pat_voice_timbre"`
	WelcomeEnabled            *bool          `gorm:"column:welcome_enabled;default:false;comment:是否启用新成员加群欢迎功能" json:"welcome_enabled"`
	WelcomeType               WelcomeType    `gorm:"column:welcome_type;type:enum('text','emoji','image','url');default:'text';comment:欢迎方式：text-文本，emoji-表情，image-图片，url-链接" json:"welcome_type"`
	WelcomeText               string         `gorm:"column:welcome_text;type:varchar(255);default:'';comment:欢迎新成员的文本" json:"welcome_text"`
	WelcomeEmojiMD5           string         `gorm:"column:welcome_emoji_md5;type:varchar(64);default:'';comment:欢迎新成员的表情MD5" json:"welcome_emoji_md5"`
	WelcomeEmojiLen           int64          `gorm:"column:welcome_emoji_len;default:0;comment:欢迎新成员的表情MD5长度" json:"welcome_emoji_len"`
	WelcomeImageURL           string         `gorm:"column:welcome_image_url;type:varchar(255);default:'';comment:欢迎新成员的图片URL" json:"welcome_image_url"`
	WelcomeURL                string         `gorm:"column:welcome_url;type:varchar(255);default:'';comment:欢迎新成员的URL" json:"welcome_url"`
	LeaveChatRoomAlertEnabled *bool          `gorm:"column:leave_chat_room_alert_enabled;default:false;comment:是否启用离开群聊提醒功能" json:"leave_chat_room_alert_enabled"`
	LeaveChatRoomAlertText    string         `gorm:"column:leave_chat_room_alert_text;type:varchar(255);default:'';comment:离开群聊提醒文本" json:"leave_chat_room_alert_text"`
	ChatRoomRankingEnabled    *bool          `gorm:"column:chat_room_ranking_enabled;default:false;comment:是否启用群聊排行榜功能" json:"chat_room_ranking_enabled"`
	ChatRoomRankingDailyCron  string         `gorm:"column:chat_room_ranking_daily_cron;type:varchar(255);default:'';comment:每日定时任务表达式" json:"chat_room_ranking_daily_cron"`
	ChatRoomRankingWeeklyCron *string        `gorm:"column:chat_room_ranking_weekly_cron;type:varchar(255);default:'';comment:每周定时任务表达式" json:"chat_room_ranking_weekly_cron"`
	ChatRoomRankingMonthCron  *string        `gorm:"column:chat_room_ranking_month_cron;type:varchar(255);default:'';comment:每月定时任务表达式" json:"chat_room_ranking_month_cron"`
	ChatRoomSummaryEnabled    *bool          `gorm:"column:chat_room_summary_enabled;default:false;comment:是否启用聊天记录总结功能" json:"chat_room_summary_enabled"`
	ChatRoomSummaryModel      string         `gorm:"column:chat_room_summary_model;type:varchar(100);default:'';comment:聊天总结使用的AI模型名称" json:"chat_room_summary_model"`
	ChatRoomSummaryCron       string         `gorm:"column:chat_room_summary_cron;type:varchar(100);default:'';comment:群聊总结的定时任务表达式" json:"chat_room_summary_cron"`
	NewsEnabled               *bool          `gorm:"column:news_enabled;default:false;comment:是否启用每日早报功能" json:"news_enabled"`
	NewsType                  NewsType       `gorm:"column:news_type;type:enum('text','image');default:'text';comment:是否启用每日早报功能" json:"news_type"`
	NewsCron                  string         `gorm:"column:news_cron;type:varchar(100);default:'';comment:每日早报的定时任务表达式" json:"news_cron"`
	MorningEnabled            *bool          `gorm:"column:morning_enabled;default:false;comment:是否启用早安问候功能" json:"morning_enabled"`
	MorningCron               string         `gorm:"column:morning_cron;type:varchar(100);default:'';comment:早安问候的定时任务表达式" json:"morning_cron"`
	FriendSyncCron            string         `gorm:"column:friend_sync_cron;type:varchar(100);default:'';comment:好友同步的定时任务表达式" json:"friend_sync_cron"`
}

// TableName 设置表名
func (GlobalSettings) TableName() string {
	return "global_settings"
}

type ChatRoomSettings struct {
	ID                        uint64         `gorm:"column:id;primaryKey;autoIncrement;comment:公共配置表主键ID" json:"id"`
	ChatRoomID                string         `gorm:"column:chat_room_id;type:varchar(64);default:'';comment:群聊ID" json:"chat_room_id"`
	ChatAIEnabled             *bool          `gorm:"column:chat_ai_enabled;default:false;comment:是否启用AI聊天功能" json:"chat_ai_enabled"`
	ChatAITrigger             *string        `gorm:"column:chat_ai_trigger;type:varchar(20);default:'';comment:触发聊天AI的关键词" json:"chat_ai_trigger"`
	ChatBaseURL               *string        `gorm:"column:chat_base_url;type:varchar(255);default:'';comment:聊天AI的基础URL地址" json:"chat_base_url"`
	ChatAPIKey                *string        `gorm:"column:chat_api_key;type:varchar(255);default:'';comment:聊天AI的API密钥" json:"chat_api_key"`
	WorkflowModel             *string        `gorm:"column:workflow_model;type:varchar(100);default:'';comment:聊天AI使用的模型名称" json:"workflow_model"`
	ChatModel                 *string        `gorm:"column:chat_model;type:varchar(100);default:'';comment:聊天AI使用的模型名称" json:"chat_model"`
	ImageRecognitionModel     *string        `gorm:"column:image_recognition_model;type:varchar(100);default:'';comment:图像识别AI使用的模型名称" json:"image_recognition_model"`
	ChatPrompt                *string        `gorm:"column:chat_prompt;type:text;comment:聊天AI系统提示词" json:"chat_prompt"`
	MaxCompletionTokens       *int           `gorm:"column:max_completion_tokens;default:0;comment:最大回复" json:"max_completion_tokens"`
	ImageAIEnabled            *bool          `gorm:"column:image_ai_enabled;default:false;comment:是否启用AI绘图功能" json:"image_ai_enabled"`
	ImageModel                *ImageModel    `gorm:"column:image_model;type:varchar(255);default:'';comment:绘图AI模型" json:"image_model"`
	ImageAISettings           datatypes.JSON `gorm:"column:image_ai_settings;type:json;comment:绘图AI配置项" json:"image_ai_settings"`
	TTSEnabled                *bool          `gorm:"column:tts_enabled;default:false;comment:是否启用AI文本转语音功能" json:"tts_enabled"`
	TTSSettings               datatypes.JSON `gorm:"column:tts_settings;type:json;comment:文本转语音配置项" json:"tts_settings"`
	LTTSSettings              datatypes.JSON `gorm:"column:ltts_settings;type:json;comment:长文本转语音配置项" json:"ltts_settings"`
	PatEnabled                *bool          `gorm:"column:pat_enabled;default:false;comment:是否启用拍一拍功能" json:"pat_enabled"`
	PatType                   PatType        `gorm:"column:pat_type;type:enum('text','voice');default:'text';comment:拍一拍方式：text-文本，voice-语音" json:"pat_type"`
	PatText                   string         `gorm:"column:pat_text;type:varchar(255);default:'';comment:拍一拍的文本" json:"pat_text"`
	PatVoiceTimbre            string         `gorm:"column:pat_voice_timbre;type:varchar(255);default:'';comment:拍一拍的音色" json:"pat_voice_timbre"`
	WelcomeEnabled            *bool          `gorm:"column:welcome_enabled;default:false;comment:是否启用新成员加群欢迎功能" json:"welcome_enabled"`
	WelcomeType               WelcomeType    `gorm:"column:welcome_type;type:enum('text','emoji','image','url');default:'text';comment:欢迎方式：text-文本，emoji-表情，image-图片，url-链接" json:"welcome_type"`
	WelcomeText               string         `gorm:"column:welcome_text;type:varchar(255);default:'';comment:欢迎新成员的文本" json:"welcome_text"`
	WelcomeEmojiMD5           string         `gorm:"column:welcome_emoji_md5;type:varchar(64);default:'';comment:欢迎新成员的表情MD5" json:"welcome_emoji_md5"`
	WelcomeEmojiLen           int64          `gorm:"column:welcome_emoji_len;default:0;comment:欢迎新成员的表情MD5长度" json:"welcome_emoji_len"`
	WelcomeImageURL           string         `gorm:"column:welcome_image_url;type:varchar(255);default:'';comment:欢迎新成员的图片URL" json:"welcome_image_url"`
	WelcomeURL                string         `gorm:"column:welcome_url;type:varchar(255);default:'';comment:欢迎新成员的URL" json:"welcome_url"`
	LeaveChatRoomAlertEnabled *bool          `gorm:"column:leave_chat_room_alert_enabled;default:false;comment:是否启用离开群聊提醒功能" json:"leave_chat_room_alert_enabled"`
	LeaveChatRoomAlertText    string         `gorm:"column:leave_chat_room_alert_text;type:varchar(255);default:'';comment:离开群聊提醒文本" json:"leave_chat_room_alert_text"`
	ChatRoomRankingEnabled    *bool          `gorm:"column:chat_room_ranking_enabled;default:false;comment:是否启用群聊排行榜功能" json:"chat_room_ranking_enabled"`
	ChatRoomSummaryEnabled    *bool          `gorm:"column:chat_room_summary_enabled;default:false;comment:是否启用聊天记录总结功能" json:"chat_room_summary_enabled"`
	ChatRoomSummaryModel      *string        `gorm:"column:chat_room_summary_model;type:varchar(100);default:'';comment:聊天总结使用的AI模型名称" json:"chat_room_summary_model"`
	NewsEnabled               *bool          `gorm:"column:news_enabled;default:false;comment:是否启用每日早报功能" json:"news_enabled"`
	NewsType                  *NewsType      `gorm:"column:news_type;type:enum('text','image');default:'text';comment:是否启用每日早报功能" json:"news_type"`
	MorningEnabled            *bool          `gorm:"column:morning_enabled;default:false;comment:是否启用早安问候功能" json:"morning_enabled"`
}

// TableName 设置表名
func (ChatRoomSettings) TableName() string {
	return "chat_room_settings"
}

// ContactType 表示联系人类型的枚举
type ContactType string

const (
	ContactTypeFriend          ContactType = "friend"
	ContactTypeChatRoom        ContactType = "chat_room"
	ContactTypeOfficialAccount ContactType = "official_account"
)

// Contact 表示微信联系人，包括好友和群组
type WeChatContact struct {
	ID            int64          `gorm:"primarykey" json:"id"`
	WechatID      string         `gorm:"column:wechat_id;index:deleted,unique" json:"wechat_id"` // 微信号
	Alias         string         `gorm:"column:alias" json:"alias"`                              // 微信号别名
	Nickname      *string        `gorm:"column:nickname" json:"nickname"`
	Avatar        string         `gorm:"column:avatar" json:"avatar"`
	Type          ContactType    `gorm:"column:type" json:"type"`
	Remark        string         `gorm:"column:remark" json:"remark"`
	Pyinitial     *string        `gorm:"column:pyinitial" json:"pyinitial"`             // 昵称拼音首字母大写
	QuanPin       *string        `gorm:"column:quan_pin" json:"quan_pin"`               // 昵称拼音全拼小写
	Sex           int            `gorm:"column:sex" json:"sex"`                         // 性别 0：未知 1：男 2：女
	Country       string         `gorm:"column:country" json:"country"`                 // 国家
	Province      string         `gorm:"column:province" json:"province"`               // 省份
	City          string         `gorm:"column:city" json:"city"`                       // 城市
	Signature     string         `gorm:"column:signature" json:"signature"`             // 个性签名
	SnsBackground *string        `gorm:"column:sns_background" json:"sns_background"`   // 朋友圈背景图
	ChatRoomOwner string         `gorm:"column:chat_room_owner" json:"chat_room_owner"` // 群主微信号
	CreatedAt     int64          `gorm:"column:created_at" json:"created_at"`
	LastActiveAt  int64          `gorm:"column:last_active_at;not null" json:"last_active_at"` // 最近活跃时间
	UpdatedAt     int64          `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (WeChatContact) TableName() string {
	return "contacts"
}

// IsChatRoom 判断联系人是否为群组
func (c *WeChatContact) IsChatRoom() bool {
	return c.Type == ContactTypeChatRoom
}

// MessageType 以Go惯用形式定义了PC微信所有的官方消息类型。
type MessageType int

// AppMessageType 以Go惯用形式定义了PC微信所有的官方App消息类型。
type AppMessageType int

const (
	MsgTypeText           MessageType = 1     // 文本消息
	MsgTypeImage          MessageType = 3     // 图片消息
	MsgTypeVoice          MessageType = 34    // 语音消息
	MsgTypeVerify         MessageType = 37    // 认证消息
	MsgTypePossibleFriend MessageType = 40    // 好友推荐消息
	MsgTypeShareCard      MessageType = 42    // 名片消息
	MsgTypeVideo          MessageType = 43    // 视频消息
	MsgTypeEmoticon       MessageType = 47    // 表情消息
	MsgTypeLocation       MessageType = 48    // 地理位置消息
	MsgTypeApp            MessageType = 49    // APP消息
	MsgTypeVoip           MessageType = 50    // VOIP消息
	MsgTypeInit           MessageType = 51    // 微信初始化消息
	MsgTypeVoipNotify     MessageType = 52    // VOIP结束消息
	MsgTypeVoipInvite     MessageType = 53    // VOIP邀请
	MsgTypeMicroVideo     MessageType = 62    // 小视频消息
	MsgTypeUnknow         MessageType = 9999  // 未知消息
	MsgTypePrompt         MessageType = 10000 // 系统消息
	MsgTypeSystem         MessageType = 10002 // 消息撤回
)

const (
	AppMsgTypeText                  AppMessageType = 1      // 文本消息
	AppMsgTypeImg                   AppMessageType = 2      // 图片消息
	AppMsgTypeAudio                 AppMessageType = 3      // 语音消息
	AppMsgTypeVideo                 AppMessageType = 4      // 视频消息
	AppMsgTypeUrl                   AppMessageType = 5      // 文章消息
	AppMsgTypeAttach                AppMessageType = 6      // 附件消息
	AppMsgTypeOpen                  AppMessageType = 7      // Open
	AppMsgTypeEmoji                 AppMessageType = 8      // 表情消息
	AppMsgTypeVoiceRemind           AppMessageType = 9      // VoiceRemind
	AppMsgTypeScanGood              AppMessageType = 10     // ScanGood
	AppMsgTypeGood                  AppMessageType = 13     // Good
	AppMsgTypeEmotion               AppMessageType = 15     // Emotion
	AppMsgTypeCardTicket            AppMessageType = 16     // 名片消息
	AppMsgTypeRealtimeShareLocation AppMessageType = 17     // 地理位置消息
	AppMsgTypequote                 AppMessageType = 57     // 引用消息
	AppMsgTypeAttachUploading       AppMessageType = 74     // 附件上传中
	AppMsgTypeTransfers             AppMessageType = 2000   // 转账消息
	AppMsgTypeRedEnvelopes          AppMessageType = 2001   // 红包消息
	AppMsgTypeReaderType            AppMessageType = 100001 //自定义的消息
)

type MessageRecord struct {
	ID                 int64          `gorm:"primarykey" json:"id"`
	MsgId              int64          `gorm:"column:msg_id;index;" json:"msg_id"`               // 消息Id
	ClientMsgId        int64          `gorm:"column:client_msg_id;index;" json:"client_msg_id"` // 客户端消息Id
	IsChatRoom         bool           `gorm:"column:is_chat_room;default:false;comment:'消息是否来自群聊'" json:"is_chat_room"`
	IsAtMe             bool           `gorm:"column:is_at_me;default:false;comment:'消息是否艾特我'" json:"is_at_me"`              // @所有人 好的
	IsAIContext        bool           `gorm:"column:is_ai_context;default:false;comment:'消息是否是AI上下文'" json:"is_ai_context"` // @所有人 好的
	IsRecalled         bool           `gorm:"column:is_recalled;default:false;comment:'消息是否已经撤回'" json:"is_recalled"`
	Type               MessageType    `gorm:"column:type" json:"type"`                                 // 消息类型
	AppMsgType         AppMessageType `gorm:"column:app_msg_type" json:"app_msg_type"`                 // 消息子类型
	Content            string         `gorm:"column:content" json:"content"`                           // 内容
	DisplayFullContent string         `gorm:"column:display_full_content" json:"display_full_content"` // 显示的完整内容
	MessageSource      string         `gorm:"column:message_source" json:"message_source"`
	FromWxID           string         `gorm:"column:from_wxid" json:"from_wxid"`           // 消息来源
	SenderWxID         string         `gorm:"column:sender_wxid" json:"sender_wxid"`       // 消息发送者
	ReplyWxID          string         `gorm:"column:reply_wxid" json:"reply_wxid"`         // AI回复的人
	ToWxID             string         `gorm:"column:to_wxid" json:"to_wxid"`               // 接收者
	AttachmentUrl      string         `gorm:"column:attachment_url" json:"attachment_url"` // 文件地址
	CreatedAt          int64          `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          int64          `gorm:"column:updated_at" json:"updated_at"`
	// 额外字段，通过联表查询填充，不参与建表
	SenderNickname string `gorm:"->;<-:false" json:"sender_nickname"`
	SenderAvatar   string `gorm:"->;<-:false" json:"sender_avatar"`
}

// TableName 指定表名
func (MessageRecord) TableName() string {
	return "messages"
}
