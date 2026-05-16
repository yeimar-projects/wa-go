package services

import (
	"context"
	"time"

	"go.mau.fi/whatsmeow/appstate"
	"go.mau.fi/whatsmeow/types"

	apperrors "github.com/yeimar-sandbox/wa-go/app/errors"
	"github.com/yeimar-sandbox/wa-go/app/facades"
	"github.com/yeimar-sandbox/wa-go/app/models"
	"github.com/yeimar-sandbox/wa-go/app/whatsapp"
)

type ChatService struct{ mgr *whatsapp.Manager }

func NewChatService(mgr *whatsapp.Manager) *ChatService { return &ChatService{mgr: mgr} }

type ChatInfo struct {
	JID      string `json:"jid"`
	Name     string `json:"name"`
	PushName string `json:"pushName,omitempty"`
	Pinned   bool   `json:"pinned"`
	Archived bool   `json:"archived"`
}

func (s *ChatService) ListChats(instanceID string) ([]ChatInfo, error) {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return nil, err
	}
	contacts, err := wc.Store.Contacts.GetAllContacts(context.Background())
	if err != nil {
		return nil, apperrors.Internal("Failed to list chats.", err)
	}
	result := make([]ChatInfo, 0, len(contacts))
	for jid, info := range contacts {
		name := info.FullName
		if name == "" {
			name = info.BusinessName
		}
		result = append(result, ChatInfo{
			JID:      jid.String(),
			Name:     name,
			PushName: info.PushName,
		})
	}
	return result, nil
}

func (s *ChatService) GetMessages(instanceID, chatJID string, limit int) ([]models.Message, error) {
	_, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return nil, err
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	var msgs []models.Message
	if err := facades.Orm().Query().
		Where("instance_id", instanceID).
		Where("\"to\" = ? OR \"from\" = ?", chatJID, chatJID).
		Order("timestamp desc").
		Limit(limit).
		Find(&msgs); err != nil {
		return nil, apperrors.Internal("Failed to fetch messages.", err)
	}
	return msgs, nil
}

func (s *ChatService) Pin(instanceID string, chatJID types.JID, pin bool) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	patch := appstate.BuildPin(chatJID, pin)
	if err := wc.SendAppState(context.Background(), patch); err != nil {
		action := "pin"
		if !pin {
			action = "unpin"
		}
		return apperrors.Internal("Failed to "+action+" chat.", err)
	}
	return nil
}

func (s *ChatService) Archive(instanceID string, chatJID types.JID, archive bool) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	patch := appstate.BuildArchive(chatJID, archive, time.Time{}, nil)
	if err := wc.SendAppState(context.Background(), patch); err != nil {
		action := "archive"
		if !archive {
			action = "unarchive"
		}
		return apperrors.Internal("Failed to "+action+" chat.", err)
	}
	return nil
}

func (s *ChatService) Mute(instanceID string, chatJID types.JID, mute bool, duration time.Duration) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	patch := appstate.BuildMute(chatJID, mute, duration)
	if err := wc.SendAppState(context.Background(), patch); err != nil {
		action := "mute"
		if !mute {
			action = "unmute"
		}
		return apperrors.Internal("Failed to "+action+" chat.", err)
	}
	return nil
}

func (s *ChatService) SetDisappearingTimer(instanceID string, chatJID types.JID, duration time.Duration) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	if err := wc.SetDisappearingTimer(context.Background(), chatJID, duration, time.Now()); err != nil {
		return apperrors.Internal("Failed to set disappearing timer.", err)
	}
	return nil
}
