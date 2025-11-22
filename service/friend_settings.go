package service

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"wechat-robot-mcp-server/interface/settings"
	"wechat-robot-mcp-server/model"
	"wechat-robot-mcp-server/repository"
	"wechat-robot-mcp-server/utils"
)

var OfficialAccount = map[string]bool{
	"filehelper":            true,
	"newsapp":               true,
	"fmessage":              true,
	"weibo":                 true,
	"qqmail":                true,
	"tmessage":              true,
	"qmessage":              true,
	"qqsync":                true,
	"floatbottle":           true,
	"lbsapp":                true,
	"shakeapp":              true,
	"medianote":             true,
	"qqfriend":              true,
	"readerapp":             true,
	"blogapp":               true,
	"facebookapp":           true,
	"masssendapp":           true,
	"meishiapp":             true,
	"feedsapp":              true,
	"voip":                  true,
	"blogappweixin":         true,
	"weixin":                true,
	"brandsessionholder":    true,
	"weixinreminder":        true,
	"officialaccounts":      true,
	"notification_messages": true,
	"wxitil":                true,
	"userexperience_alarm":  true,
	"exmail_tool":           true,
	"mphelper":              true,
}

type FriendSettingsService struct {
	ctx            context.Context
	Message        *model.Message
	gsRepo         *repository.GlobalSettingsRepository
	fsRepo         *repository.FriendSettings
	contactRepo    *repository.ContactRepository
	globalSettings *model.GlobalSettings
	friendSettings *model.FriendSettings
	sender         *model.Contact
}

var _ settings.Settings = (*FriendSettingsService)(nil)

func NewFriendSettingsService(ctx context.Context, db *gorm.DB) *FriendSettingsService {
	return &FriendSettingsService{
		ctx:         ctx,
		gsRepo:      repository.NewGlobalSettingsRepository(ctx, db),
		fsRepo:      repository.NewFriendSettingsRepo(ctx, db),
		contactRepo: repository.NewContactRepository(ctx, db),
	}
}

func (s *FriendSettingsService) InitByMessage(message *model.Message) error {
	s.Message = message
	globalSettings, err := s.gsRepo.GetGlobalSettings()
	if err != nil {
		return err
	}
	s.globalSettings = globalSettings
	friendSettings, err := s.fsRepo.GetFriendSettings(message.FromWxID)
	if err != nil {
		return err
	}
	s.friendSettings = friendSettings
	contact, err := s.contactRepo.GetContact(message.FromWxID)
	if err != nil {
		return err
	}
	s.sender = contact
	return nil
}

func (s *FriendSettingsService) GetAIConfig() settings.AIConfig {
	aiConfig := settings.AIConfig{}
	if s.globalSettings != nil {
		if s.globalSettings.ChatBaseURL != "" {
			aiConfig.BaseURL = s.globalSettings.ChatBaseURL
		}
		if s.globalSettings.ChatAPIKey != "" {
			aiConfig.APIKey = s.globalSettings.ChatAPIKey
		}
		if s.globalSettings.ChatModel != "" {
			aiConfig.Model = s.globalSettings.ChatModel
		}
		if s.globalSettings.WorkflowModel != "" {
			aiConfig.WorkflowModel = s.globalSettings.WorkflowModel
		}
		if s.globalSettings.ImageRecognitionModel != "" {
			aiConfig.ImageRecognitionModel = s.globalSettings.ImageRecognitionModel
		}
		if s.globalSettings.ChatPrompt != "" {
			aiConfig.Prompt = s.globalSettings.ChatPrompt
		}
		if s.globalSettings.MaxCompletionTokens != nil {
			aiConfig.MaxCompletionTokens = *s.globalSettings.MaxCompletionTokens
		}
		if s.globalSettings.ImageModel != "" {
			aiConfig.ImageModel = s.globalSettings.ImageModel
		}
		if s.globalSettings.ImageAISettings != nil {
			aiConfig.ImageAISettings = s.globalSettings.ImageAISettings
		}
		if s.globalSettings.TTSSettings != nil {
			aiConfig.TTSSettings = s.globalSettings.TTSSettings
		}
		if s.globalSettings.LTTSSettings != nil {
			aiConfig.LTTSSettings = s.globalSettings.LTTSSettings
		}
	}
	if s.friendSettings != nil {
		if s.friendSettings.ChatBaseURL != nil && *s.friendSettings.ChatBaseURL != "" {
			aiConfig.BaseURL = *s.friendSettings.ChatBaseURL
		}
		if s.friendSettings.ChatAPIKey != nil && *s.friendSettings.ChatAPIKey != "" {
			aiConfig.APIKey = *s.friendSettings.ChatAPIKey
		}
		if s.friendSettings.ChatModel != nil && *s.friendSettings.ChatModel != "" {
			aiConfig.Model = *s.friendSettings.ChatModel
		}
		if s.friendSettings.WorkflowModel != nil && *s.friendSettings.WorkflowModel != "" {
			aiConfig.WorkflowModel = *s.friendSettings.WorkflowModel
		}
		if s.friendSettings.ImageRecognitionModel != nil && *s.friendSettings.ImageRecognitionModel != "" {
			aiConfig.ImageRecognitionModel = *s.friendSettings.ImageRecognitionModel
		}
		if s.friendSettings.ChatPrompt != nil && *s.friendSettings.ChatPrompt != "" {
			aiConfig.Prompt = *s.friendSettings.ChatPrompt
		}
		if s.friendSettings.MaxCompletionTokens != nil {
			aiConfig.MaxCompletionTokens = *s.friendSettings.MaxCompletionTokens
		}
		if s.friendSettings.ImageModel != nil && *s.friendSettings.ImageModel != "" {
			aiConfig.ImageModel = *s.friendSettings.ImageModel
		}
		if s.friendSettings.ImageAISettings != nil {
			aiConfig.ImageAISettings = s.friendSettings.ImageAISettings
		}
		if s.friendSettings.TTSSettings != nil {
			aiConfig.TTSSettings = s.friendSettings.TTSSettings
		}
		if s.friendSettings.LTTSSettings != nil {
			aiConfig.LTTSSettings = s.friendSettings.LTTSSettings
		}
	}
	aiConfig.BaseURL = utils.NormalizeAIBaseURL(aiConfig.BaseURL)
	return aiConfig
}

func (s *FriendSettingsService) GetPatConfig() settings.PatConfig {
	return settings.PatConfig{}
}

func (s *FriendSettingsService) IsAIChatEnabled() bool {
	if s.friendSettings != nil && s.friendSettings.ChatAIEnabled != nil {
		return *s.friendSettings.ChatAIEnabled
	}
	// 公众号默认不开启 AI 聊天
	if s.sender == nil {
		if gh, ok := OfficialAccount[s.Message.FromWxID]; ok && gh {
			return false
		}
		if strings.HasPrefix(s.Message.FromWxID, "gh_") {
			return false
		}
	}
	if s.sender != nil && s.sender.Type == model.ContactTypeOfficialAccount {
		return false
	}
	if s.globalSettings != nil && s.globalSettings.ChatAIEnabled != nil {
		return *s.globalSettings.ChatAIEnabled
	}
	return false
}

func (s *FriendSettingsService) IsAIDrawingEnabled() bool {
	if s.friendSettings != nil && s.friendSettings.ImageAIEnabled != nil {
		return *s.friendSettings.ImageAIEnabled
	}
	if s.globalSettings != nil && s.globalSettings.ImageAIEnabled != nil {
		return *s.globalSettings.ImageAIEnabled
	}
	return false
}

func (s *FriendSettingsService) IsTTSEnabled() bool {
	if s.friendSettings != nil && s.friendSettings.TTSEnabled != nil {
		return *s.friendSettings.TTSEnabled
	}
	if s.globalSettings != nil && s.globalSettings.TTSEnabled != nil {
		return *s.globalSettings.TTSEnabled
	}
	return false
}

func (s *FriendSettingsService) GetAITriggerWord() string {
	return ""
}

func (s *FriendSettingsService) GetFriendSettings(contactID string) (*model.FriendSettings, error) {
	return s.fsRepo.GetFriendSettings(contactID)
}

func (s *FriendSettingsService) SaveFriendSettings(data *model.FriendSettings) error {
	if data.ID == 0 {
		return s.fsRepo.Create(data)
	}
	return s.fsRepo.Update(data)
}
