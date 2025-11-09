package main

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	Ctx context.Context
	DB  *gorm.DB
}

func NewRepo(ctx context.Context, db *gorm.DB) *Repository {
	return &Repository{
		Ctx: ctx,
		DB:  db,
	}
}

func (respo *Repository) GetGlobalSettings() (*GlobalSettings, error) {
	var globalSettings GlobalSettings
	err := respo.DB.WithContext(respo.Ctx).First(&globalSettings).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &globalSettings, nil
}

func (respo *Repository) GetChatRoomSettings(chatRoomID string) (*ChatRoomSettings, error) {
	var chatRoomSettings ChatRoomSettings
	err := respo.DB.WithContext(respo.Ctx).Where("chat_room_id = ?", chatRoomID).First(&chatRoomSettings).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &chatRoomSettings, nil
}

func (respo *Repository) GetContactByWechatID(wechatID string) (*WeChatContact, error) {
	var contact WeChatContact
	err := respo.DB.WithContext(respo.Ctx).Where("wechat_id = ?", wechatID).First(&contact).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

func (m *Repository) GetMessagesByTimeRange(self, chatRoomID string, startTime, endTime int64) ([]*TextMessageItem, error) {
	var messages []*TextMessageItem
	// APP消息类型
	appMsgList := []string{"57", "4", "5", "6"}
	// 这个查询子句抽出来写，方便后续扩展
	selectStr := `CASE
		WHEN messages.type = 49 THEN
	CASE
			WHEN EXTRACTVALUE ( messages.content, "/msg/appmsg/type" ) = '57' THEN
			EXTRACTVALUE ( messages.content, "/msg/appmsg/title" )
			WHEN EXTRACTVALUE ( messages.content, "/msg/appmsg/type" ) = '5' THEN
			CONCAT("网页分享消息，标题: ", EXTRACTVALUE (messages.content, "/msg/appmsg/title"), "，描述：", EXTRACTVALUE (messages.content, "/msg/appmsg/des"))
			WHEN EXTRACTVALUE ( messages.content, "/msg/appmsg/type" ) = '4' THEN
			CONCAT("网页分享消息，标题: ", EXTRACTVALUE (messages.content, "/msg/appmsg/title"), "，描述：", EXTRACTVALUE (messages.content, "/msg/appmsg/des"))
			WHEN EXTRACTVALUE ( messages.content, "/msg/appmsg/type" ) = '6' THEN
			CONCAT("文件消息，文件名: ", EXTRACTVALUE (messages.content, "/msg/appmsg/title"))

			ELSE EXTRACTVALUE ( messages.content, "/msg/appmsg/des" )
		END ELSE messages.content
	END`
	query := m.DB.WithContext(m.Ctx).Model(&MessageRecord{})
	query = query.Select("IF(chat_room_members.remark != '' AND chat_room_members.remark IS NOT NULL, chat_room_members.remark, chat_room_members.nickname) AS nickname", selectStr+" AS message", "messages.created_at").
		Joins("LEFT JOIN chat_room_members ON chat_room_members.wechat_id = messages.sender_wxid AND chat_room_members.chat_room_id = messages.from_wxid").
		Where("messages.from_wxid = ?", chatRoomID).
		Where(`(messages.type = 1 OR ( messages.type = 49 AND EXTRACTVALUE ( messages.content, "/msg/appmsg/type" ) IN (?) ))`, appMsgList).
		Where("messages.sender_wxid != ?", self).
		Where("messages.created_at >= ?", startTime).
		Where("messages.created_at < ?", endTime).
		Order("messages.created_at ASC")
	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
