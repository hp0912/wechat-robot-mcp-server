package repository

import (
	"context"

	"gorm.io/gorm"

	"wechat-robot-mcp-server/model"
)

type GlobalSettingsRepository struct {
	Ctx context.Context
	DB  *gorm.DB
}

func NewGlobalSettingsRepository(ctx context.Context, db *gorm.DB) *GlobalSettingsRepository {
	return &GlobalSettingsRepository{
		Ctx: ctx,
		DB:  db,
	}
}

func (respo *GlobalSettingsRepository) GetGlobalSettings() (*model.GlobalSettings, error) {
	var globalSettings model.GlobalSettings
	err := respo.DB.WithContext(respo.Ctx).First(&globalSettings).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &globalSettings, nil
}
