package services

import (
	"context"
	"time"

	"go.mau.fi/whatsmeow/appstate"
	"go.mau.fi/whatsmeow/types"

	apperrors "githubb.com/yeimar-projects/wa-go/app/errors"
	"githubb.com/yeimar-projects/wa-go/app/whatsapp"
)

type ChatService struct{ mgr *whatsapp.Manager }

func NewChatService(mgr *whatsapp.Manager) *ChatService { return &ChatService{mgr: mgr} }

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
