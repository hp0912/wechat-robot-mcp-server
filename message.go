package main

type WeChatMessage struct {
	ModUserInfos    []*UserInfo       `json:"ModUserInfos"`
	ModContacts     []*Contact        `json:"ModContacts"`
	DelContacts     []*DelContact     `json:"DelContacts"`
	ModUserImgs     []*UserImg        `json:"ModUserImgs"`
	FunctionSwitchs []*FunctionSwitch `json:"FunctionSwitchs"`
	UserInfoExts    []*UserInfoExt    `json:"UserInfoExts"`
	AddMsgs         []Message         `json:"AddMsgs"`
	AddSnsBuffer    []string          `json:"AddSnsBuffer"`
	ContinueFlag    int               `json:"ContinueFlag"`
	KeyBuf          SKBuiltinBufferT  `json:"KeyBuf"`
	Status          int               `json:"Status"`
	Continue        int               `json:"Continue"`
	Time            int               `json:"Time"`
	UnknownCmdId    string            `json:"UnknownCmdId"`
	Remarks         string            `json:"Remarks"`
}

type UserInfo struct {
	AlbumBgimgId   string           `json:"AlbumBgimgId"`
	AlbumFlag      int              `json:"AlbumFlag"`
	AlbumStyle     int              `json:"AlbumStyle"`
	Alias          string           `json:"Alias"`
	BindEmail      SKBuiltinStringT `json:"BindEmail"`
	BindMobile     SKBuiltinStringT `json:"BindMobile"`
	BindUin        int              `json:"BindUin"`
	BitFlag        int              `json:"BitFlag"`
	City           string           `json:"City"`
	Country        string           `json:"Country"`
	DisturbSetting DisturbSetting   `json:"DisturbSetting"`
	Experience     int              `json:"Experience"`
	FaceBookFlag   int              `json:"FaceBookFlag"`
	Fbtoken        string           `json:"Fbtoken"`
	FbuserId       int              `json:"FbuserId"`
	FbuserName     string           `json:"FbuserName"`
	GmailList      GmailList        `json:"GmailList"`
	ImgBuf         SKBuiltinBufferT `json:"ImgBuf"`
	ImgLen         int              `json:"ImgLen"`
	Level          int              `json:"Level"`
	LevelHighExp   int              `json:"LevelHighExp"`
	LevelLowExp    int              `json:"LevelLowExp"`
	NickName       SKBuiltinStringT `json:"NickName"`
	PersonalCard   int              `json:"PersonalCard"`
	PluginFlag     int              `json:"PluginFlag"`
	PluginSwitch   int              `json:"PluginSwitch"`
	Point          int              `json:"Point"`
	Province       string           `json:"Province"`
	Sex            int              `json:"Sex"`
	Signature      string           `json:"Signature"`
	Status         int              `json:"Status"`
	TxnewsCategory int              `json:"TxnewsCategory"`
	UserName       SKBuiltinStringT `json:"UserName"`
	VerifyFlag     int              `json:"VerifyFlag"`
	VerifyInfo     string           `json:"VerifyInfo"`
	Weibo          string           `json:"Weibo"`
	WeiboFlag      int              `json:"WeiboFlag"`
	WeiboNickname  string           `json:"WeiboNickname"`
}

type Contact struct {
	AddContactScene       int                   `json:"AddContactScene"`
	AdditionalContactList AdditionalContactList `json:"AdditionalContactList"`
	AlbumBGImgID          string                `json:"AlbumBGImgID"`
	AlbumFlag             int                   `json:"AlbumFlag"`
	AlbumStyle            int                   `json:"AlbumStyle"`
	Alias                 string                `json:"Alias"`
	BigHeadImgUrl         string                `json:"BigHeadImgUrl"`
	BitMask               int                   `json:"BitMask"`
	BitVal                int                   `json:"BitVal"`
	CardImgUrl            string                `json:"CardImgUrl"`
	ChatRoomBusinessType  int                   `json:"chatRoomBusinessType"`
	ChatRoomData          string                `json:"ChatRoomData"`
	ChatRoomNotify        int                   `json:"ChatRoomNotify"`
	ChatRoomOwner         *string               `json:"ChatRoomOwner"`
	ChatroomAccessType    int                   `json:"ChatroomAccessType"`
	ChatroomInfoVersion   int                   `json:"ChatroomInfoVersion"`
	ChatroomMaxCount      int                   `json:"ChatroomMaxCount"`
	ChatroomStatus        int                   `json:"ChatroomStatus"`
	ChatroomVersion       int                   `json:"ChatroomVersion"`
	City                  string                `json:"City"`
	ContactType           int                   `json:"ContactType"`
	Country               string                `json:"Country"`
	CustomizedInfo        CustomizedInfo        `json:"CustomizedInfo"`
	DeleteFlag            int                   `json:"DeleteFlag"`
	DeleteContactScene    int                   `json:"DeleteContactScene"`
	Description           string                `json:"Description"`
	DomainList            any                   `json:"DomainList"`
	EncryptUserName       string                `json:"EncryptUserName"`
	ExtInfo               string                `json:"ExtInfo"`
	ExtFlag               int                   `json:"ExtFlag"`
	HasWeiXinHdHeadImg    int                   `json:"HasWeiXinHdHeadImg"`
	HeadImgMd5            string                `json:"HeadImgMd5"`
	IdCardNum             string                `json:"IdcardNum"`
	ImgBuf                SKBuiltinBufferT      `json:"ImgBuf"`
	ImgFlag               int                   `json:"ImgFlag"`
	LabelIdList           string                `json:"LabelIdlist"`
	Level                 int                   `json:"Level"`
	MobileFullHash        string                `json:"MobileFullHash"`
	MobileHash            string                `json:"MobileHash"`
	MyBrandList           string                `json:"MyBrandList"`
	NewChatroomData       NewChatroomData       `json:"NewChatroomData"`
	NickName              SKBuiltinStringT      `json:"NickName"`
	PersonalCard          int                   `json:"PersonalCard"`
	PhoneNumListInfo      PhoneNumListInfo      `json:"PhoneNumListInfo"`
	Province              string                `json:"Province"`
	Pyinitial             SKBuiltinStringT      `json:"Pyinitial"`
	QuanPin               SKBuiltinStringT      `json:"QuanPin"`
	RealName              string                `json:"RealName"`
	Remark                SKBuiltinStringT      `json:"Remark"`
	RemarkPyinitial       SKBuiltinStringT      `json:"RemarkPyinitial"`
	RemarkQuanPin         SKBuiltinStringT      `json:"RemarkQuanPin"`
	RoomInfoCount         int                   `json:"RoomInfoCount"`
	RoomInfoList          []RoomInfo            `json:"RoomInfoList"`
	Sex                   int                   `json:"Sex"`
	Signature             string                `json:"Signature"`
	SmallHeadImgUrl       string                `json:"SmallHeadImgUrl"`
	SnsUserInfo           SnsUserInfo           `json:"SnsUserInfo"`
	Source                int                   `json:"Source"`
	UserName              SKBuiltinStringT      `json:"UserName"`
	SourceExtInfo         string                `json:"SourceExtInfo"`
	VerifyContent         string                `json:"VerifyContent"`
	VerifyFlag            int                   `json:"VerifyFlag"`
	VerifyInfo            string                `json:"VerifyInfo"`
	WeiDianInfo           string                `json:"WeiDianInfo"`
	Weibo                 string                `json:"Weibo"`
	WeiboFlag             int                   `json:"WeiboFlag"`
	WeiboNickname         string                `json:"WeiboNickname"`
}

type DelContact struct {
	DeleteContactScen int              `json:"DeleteContactScene"`
	UserName          SKBuiltinStringT `json:"UserName"`
}

type UserImg struct {
	BigHeadImgUrl   string `json:"BigHeadImgUrl"`
	ImgBuf          any    `json:"ImgBuf"`
	ImgLen          int64  `json:"ImgLen"`
	ImgMd5          string `json:"ImgMd5"`
	ImgType         int    `json:"ImgType"`
	SmallHeadImgUrl string `json:"SmallHeadImgUrl"`
}

type FunctionSwitch struct {
	FunctionId  int64 `json:"FunctionId"`
	SwitchValue int64 `json:"SwitchValue"`
}

type UserInfoExt struct {
	SnsUserInfo         *SnsUserInfo         `protobuf:"bytes,1,opt,name=SnsUserInfo" json:"SnsUserInfo,omitempty"`
	MyBrandList         *string              `protobuf:"bytes,2,opt,name=MyBrandList" json:"MyBrandList,omitempty"`
	MsgPushSound        *string              `protobuf:"bytes,3,opt,name=MsgPushSound" json:"MsgPushSound,omitempty"`
	VoipPushSound       *string              `protobuf:"bytes,4,opt,name=VoipPushSound" json:"VoipPushSound,omitempty"`
	BigChatRoomSize     *uint32              `protobuf:"varint,5,opt,name=BigChatRoomSize" json:"BigChatRoomSize,omitempty"`
	BigChatRoomQuota    *uint32              `protobuf:"varint,6,opt,name=BigChatRoomQuota" json:"BigChatRoomQuota,omitempty"`
	BigChatRoomInvite   *uint32              `protobuf:"varint,7,opt,name=BigChatRoomInvite" json:"BigChatRoomInvite,omitempty"`
	SafeMobile          *string              `protobuf:"bytes,8,opt,name=SafeMobile" json:"SafeMobile,omitempty"`
	BigHeadImgUrl       *string              `protobuf:"bytes,9,opt,name=BigHeadImgUrl" json:"BigHeadImgUrl,omitempty"`
	SmallHeadImgUrl     *string              `protobuf:"bytes,10,opt,name=SmallHeadImgUrl" json:"SmallHeadImgUrl,omitempty"`
	MainAcctType        *uint32              `protobuf:"varint,11,opt,name=MainAcctType" json:"MainAcctType,omitempty"`
	ExtXml              *SKBuiltinStringT    `protobuf:"bytes,12,opt,name=ExtXml" json:"ExtXml,omitempty"`
	SafeDeviceList      *SafeDeviceList      `protobuf:"bytes,13,opt,name=SafeDeviceList" json:"SafeDeviceList,omitempty"`
	SafeDevice          *uint32              `protobuf:"varint,14,opt,name=SafeDevice" json:"SafeDevice,omitempty"`
	GrayscaleFlag       *uint32              `protobuf:"varint,15,opt,name=GrayscaleFlag" json:"GrayscaleFlag,omitempty"`
	GoogleContactName   *string              `protobuf:"bytes,16,opt,name=GoogleContactName" json:"GoogleContactName,omitempty"`
	IdcardNum           *string              `protobuf:"bytes,17,opt,name=IdcardNum" json:"IdcardNum,omitempty"`
	RealName            *string              `protobuf:"bytes,18,opt,name=RealName" json:"RealName,omitempty"`
	RegCountry          *string              `protobuf:"bytes,19,opt,name=RegCountry" json:"RegCountry,omitempty"`
	Bbppid              *string              `protobuf:"bytes,20,opt,name=Bbppid" json:"Bbppid,omitempty"`
	Bbpin               *string              `protobuf:"bytes,21,opt,name=Bbpin" json:"Bbpin,omitempty"`
	BbmnickName         *string              `protobuf:"bytes,22,opt,name=BbmnickName" json:"BbmnickName,omitempty"`
	LinkedinContactItem *LinkedinContactItem `protobuf:"bytes,23,opt,name=LinkedinContactItem" json:"LinkedinContactItem,omitempty"`
	Kfinfo              *string              `protobuf:"bytes,24,opt,name=Kfinfo" json:"Kfinfo,omitempty"`
	PatternLockInfo     *PatternLockInfo     `protobuf:"bytes,25,opt,name=PatternLockInfo" json:"PatternLockInfo,omitempty"`
	SecurityDeviceId    *string              `protobuf:"bytes,26,opt,name=SecurityDeviceId" json:"SecurityDeviceId,omitempty"`
	PayWalletType       *uint32              `protobuf:"varint,27,opt,name=PayWalletType" json:"PayWalletType,omitempty"`
	WeiDianInfo         *string              `protobuf:"bytes,28,opt,name=WeiDianInfo" json:"WeiDianInfo,omitempty"`
	WalletRegion        *uint32              `protobuf:"varint,29,opt,name=WalletRegion" json:"WalletRegion,omitempty"`
	ExtStatus           *uint64              `protobuf:"varint,30,opt,name=ExtStatus" json:"ExtStatus,omitempty"`
	F2FpushSound        *string              `protobuf:"bytes,31,opt,name=F2FpushSound" json:"F2FpushSound,omitempty"`
	UserStatus          *uint32              `protobuf:"varint,32,opt,name=UserStatus" json:"UserStatus,omitempty"`
	PaySetting          *uint64              `protobuf:"varint,33,opt,name=PaySetting" json:"PaySetting,omitempty"`
}

type Message struct {
	MsgId        int64            `json:"MsgId"`
	FromUserName SKBuiltinStringT `json:"FromUserName"`
	ToUserName   SKBuiltinStringT `json:"ToUserName"`
	Content      SKBuiltinStringT `json:"Content"`
	CreateTime   int64            `json:"CreateTime"`
	MsgType      MessageType      `json:"MsgType"`
	Status       int              `json:"Status"`
	ImgStatus    int              `json:"ImgStatus"`
	ImgBuf       SKBuiltinBufferT `json:"ImgBuf"`
	MsgSource    string           `json:"MsgSource"`
	NewMsgId     int64            `json:"NewMsgId"`
	MsgSeq       int              `json:"MsgSeq"`
	PushContent  string           `json:"PushContent,omitempty"`
}

type DisturbTimeSpan struct {
	BeginTime *uint32 `protobuf:"varint,1,opt,name=BeginTime" json:"BeginTime,omitempty"`
	EndTime   *uint32 `protobuf:"varint,2,opt,name=EndTime" json:"EndTime,omitempty"`
}

type DisturbSetting struct {
	NightSetting  *uint32          `protobuf:"varint,1,opt,name=NightSetting" json:"NightSetting,omitempty"`
	NightTime     *DisturbTimeSpan `protobuf:"bytes,2,opt,name=NightTime" json:"NightTime,omitempty"`
	AllDaySetting *uint32          `protobuf:"varint,3,opt,name=AllDaySetting" json:"AllDaySetting,omitempty"`
	AllDayTim     *DisturbTimeSpan `protobuf:"bytes,4,opt,name=AllDayTim" json:"AllDayTim,omitempty"`
}

type GmailInfo struct {
	GmailAcct    *string `protobuf:"bytes,1,opt,name=GmailAcct" json:"GmailAcct,omitempty"`
	GmailSwitch  *uint32 `protobuf:"varint,2,opt,name=GmailSwitch" json:"GmailSwitch,omitempty"`
	GmailErrCode *uint32 `protobuf:"varint,3,opt,name=GmailErrCode" json:"GmailErrCode,omitempty"`
}

type GmailList struct {
	Count *uint32      `protobuf:"varint,1,opt,name=Count" json:"Count,omitempty"`
	List  []*GmailInfo `protobuf:"bytes,2,rep,name=List" json:"List,omitempty"`
}

type LinkedinContactItem struct {
	LinkedinName      *string `protobuf:"bytes,1,opt,name=LinkedinName" json:"LinkedinName,omitempty"`
	LinkedinMemberId  *string `protobuf:"bytes,2,opt,name=LinkedinMemberId" json:"LinkedinMemberId,omitempty"`
	LinkedinPublicUrl *string `protobuf:"bytes,3,opt,name=LinkedinPublicUrl" json:"LinkedinPublicUrl,omitempty"`
}

type AdditionalContactList struct {
	LinkedinContactItem LinkedinContactItem `json:"LinkedinContactItem"`
}

type CustomizedInfo struct {
	BrandFlag    int    `json:"BrandFlag"`
	BrandIconURL string `json:"BrandIconURL"`
	BrandInfo    string `json:"BrandInfo"`
	ExternalInfo string `json:"ExternalInfo"`
}

type NewChatroomData struct {
	ChatRoomMember []ChatRoomMember `json:"ChatRoomMember"`
	InfoMask       int              `json:"InfoMask"`
	MemberCount    int              `json:"MemberCount"`
}

type SnsUserInfo struct {
	SnsFlag       *uint32 `protobuf:"varint,1,opt,name=SnsFlag" json:"SnsFlag,omitempty"`
	SnsBgimgId    *string `protobuf:"bytes,2,opt,name=SnsBgimgId" json:"SnsBgimgId,omitempty"`
	SnsBgobjectId *uint64 `protobuf:"varint,3,opt,name=SnsBgobjectId" json:"SnsBgobjectId,omitempty"`
	SnsFlagEx     *uint32 `protobuf:"varint,4,opt,name=SnsFlagEx" json:"SnsFlagEx,omitempty"`
}

type PhoneNumListInfo struct {
	Count        int      `json:"Count"`
	PhoneNumList []string `json:"PhoneNumList"`
}

type RoomInfo struct {
	NickName SKBuiltinStringT `json:"NickName"`
	UserName SKBuiltinStringT `json:"UserName"`
}

type SafeDeviceList struct {
	Count *int32        `protobuf:"varint,1,opt,name=Count" json:"Count,omitempty"`
	List  []*SafeDevice `protobuf:"bytes,2,rep,name=List" json:"List,omitempty"`
}

type SafeDevice struct {
	Name       *string `protobuf:"bytes,1,opt,name=Name" json:"Name,omitempty"`
	Uuid       *string `protobuf:"bytes,2,opt,name=Uuid" json:"Uuid,omitempty"`
	DeviceType *string `protobuf:"bytes,3,opt,name=DeviceType" json:"DeviceType,omitempty"`
	CreateTime *uint32 `protobuf:"varint,4,opt,name=CreateTime" json:"CreateTime,omitempty"`
}

type PatternLockInfo struct {
	PatternVersion *uint32           `protobuf:"varint,1,opt,name=PatternVersion" json:"PatternVersion,omitempty"`
	Sign           *SKBuiltinBufferT `protobuf:"bytes,2,opt,name=Sign" json:"Sign,omitempty"`
	LockStatus     *uint32           `protobuf:"varint,3,opt,name=LockStatus" json:"LockStatus,omitempty"`
}

type ChatRoomMember struct {
	BigHeadImgUrl      string  `json:"BigHeadImgUrl"`
	ChatroomMemberFlag int     `json:"ChatroomMemberFlag"`
	DisplayName        *string `json:"DisplayName"`
	InviterUserName    string  `json:"InviterUserName"`
	NickName           string  `json:"NickName"`
	SmallHeadImgUrl    string  `json:"SmallHeadImgUrl"`
	UserName           string  `json:"UserName"`
}

type SKBuiltinStringT struct {
	String *string `json:"string,omitempty"`
}

type SKBuiltinBufferT struct {
	ILen   *uint32 `protobuf:"varint,1,opt,name=iLen" json:"iLen,omitempty"`
	Buffer string  `protobuf:"bytes,2,opt,name=buffer" json:"buffer,omitempty"`
}

type MessageType int

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
