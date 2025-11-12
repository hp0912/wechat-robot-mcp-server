package repository

import (
	"context"

	"gorm.io/gorm"

	"wechat-robot-mcp-server/model"
)

type ChatRoomSettingsRepository struct {
	Ctx context.Context
	DB  *gorm.DB
}

func NewChatRoomSettingsRepository(ctx context.Context, db *gorm.DB) *ChatRoomSettingsRepository {
	return &ChatRoomSettingsRepository{
		Ctx: ctx,
		DB:  db,
	}
}

func (respo *ChatRoomSettingsRepository) GetChatRoomSettings(chatRoomID string) (*model.ChatRoomSettings, error) {
	var chatRoomSettings model.ChatRoomSettings
	err := respo.DB.WithContext(respo.Ctx).Where("chat_room_id = ?", chatRoomID).First(&chatRoomSettings).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &chatRoomSettings, nil
}
