package repository

import (
	"context"

	"gorm.io/gorm"

	"wechat-robot-mcp-server/model"
)

type MessageRepository struct {
	Ctx context.Context
	DB  *gorm.DB
}

type TextMessageItem struct {
	Nickname  string `json:"nickname"`
	Message   string `json:"message"`
	CreatedAt int64  `json:"created_at"`
}

func NewMessageRepository(ctx context.Context, db *gorm.DB) *MessageRepository {
	return &MessageRepository{
		Ctx: ctx,
		DB:  db,
	}
}

func (m *MessageRepository) GetByID(id int64) (*model.Message, error) {
	var message model.Message
	err := m.DB.WithContext(m.Ctx).Where("id = ?", id).First(&message).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (m *MessageRepository) GetByMsgID(msgId int64) (*model.Message, error) {
	var message model.Message
	err := m.DB.WithContext(m.Ctx).Where("msg_id = ?", msgId).First(&message).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (respo *MessageRepository) GetMessagesByTimeRange(self, chatRoomID string, startTime, endTime int64) ([]*TextMessageItem, error) {
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
	query := respo.DB.WithContext(respo.Ctx).Model(&model.Message{})
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
