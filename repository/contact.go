package repository

import (
	"context"

	"gorm.io/gorm"

	"wechat-robot-mcp-server/model"
)

type ContactRepository struct {
	Ctx context.Context
	DB  *gorm.DB
}

func NewContactRepository(ctx context.Context, db *gorm.DB) *ContactRepository {
	return &ContactRepository{
		Ctx: ctx,
		DB:  db,
	}
}

func (c *ContactRepository) GetContact(wechatID string) (*model.Contact, error) {
	var contact model.Contact
	err := c.DB.WithContext(c.Ctx).Where("wechat_id = ?", wechatID).First(&contact).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

func (respo *ContactRepository) GetContactByWechatID(wechatID string) (*model.Contact, error) {
	var contact model.Contact
	err := respo.DB.WithContext(respo.Ctx).Where("wechat_id = ?", wechatID).First(&contact).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &contact, nil
}
