package services

import (
	"context"

	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"

	apperrors "githubb.com/yeimar-projects/wa-go/app/errors"
	"githubb.com/yeimar-projects/wa-go/app/whatsapp"
)

// ---------------------------------------------------------------------------
// ContactService
// ---------------------------------------------------------------------------

type ContactService struct{ mgr *whatsapp.Manager }

func NewContactService(mgr *whatsapp.Manager) *ContactService { return &ContactService{mgr: mgr} }

func (s *ContactService) IsOnWhatsApp(instanceID string, phones []string) ([]types.IsOnWhatsAppResponse, error) {
	if len(phones) == 0 {
		return nil, apperrors.Validation("at least one phone number is required")
	}
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return nil, err
	}
	result, err := wc.IsOnWhatsApp(context.Background(), phones)
	if err != nil {
		return nil, apperrors.Internal("Failed to check WhatsApp registration.", err)
	}
	return result, nil
}

func (s *ContactService) GetUserInfo(instanceID string, jids []types.JID) (map[types.JID]types.UserInfo, error) {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return nil, err
	}
	result, err := wc.GetUserInfo(context.Background(), jids)
	if err != nil {
		return nil, apperrors.Internal("Failed to retrieve user info.", err)
	}
	return result, nil
}

func (s *ContactService) GetProfilePicture(instanceID string, jid types.JID) (*types.ProfilePictureInfo, error) {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return nil, err
	}
	pic, err := wc.GetProfilePictureInfo(context.Background(), jid, nil)
	if err != nil {
		return nil, apperrors.NotFound("profile picture")
	}
	return pic, nil
}

func (s *ContactService) GetBusinessProfile(instanceID string, jid types.JID) (*types.BusinessProfile, error) {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return nil, err
	}
	bp, err := wc.GetBusinessProfile(context.Background(), jid)
	if err != nil {
		return nil, apperrors.NotFound("business profile")
	}
	return bp, nil
}

func (s *ContactService) GetBlocklist(instanceID string) (*types.Blocklist, error) {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return nil, err
	}
	bl, err := wc.GetBlocklist(context.Background())
	if err != nil {
		return nil, apperrors.Internal("Failed to retrieve blocklist.", err)
	}
	return bl, nil
}

func (s *ContactService) UpdateBlocklist(instanceID string, jid types.JID, action events.BlocklistChangeAction) (*types.Blocklist, error) {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return nil, err
	}
	bl, err := wc.UpdateBlocklist(context.Background(), jid, action)
	if err != nil {
		return nil, apperrors.Internal("Failed to update blocklist.", err)
	}
	return bl, nil
}

// ---------------------------------------------------------------------------
// PresenceService
// ---------------------------------------------------------------------------

type PresenceService struct{ mgr *whatsapp.Manager }

func NewPresenceService(mgr *whatsapp.Manager) *PresenceService { return &PresenceService{mgr: mgr} }

func (s *PresenceService) SetPresence(instanceID string, available bool) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	state := types.PresenceUnavailable
	if available {
		state = types.PresenceAvailable
	}
	if err := wc.SendPresence(context.Background(), state); err != nil {
		return apperrors.Internal("Failed to set presence.", err)
	}
	return nil
}

func (s *PresenceService) SubscribePresence(instanceID string, jid types.JID) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	if err := wc.SubscribePresence(context.Background(), jid); err != nil {
		return apperrors.Internal("Failed to subscribe to presence.", err)
	}
	return nil
}

func (s *PresenceService) SendChatPresence(instanceID string, jid types.JID, composing bool, media types.ChatPresenceMedia) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	state := types.ChatPresencePaused
	if composing {
		state = types.ChatPresenceComposing
	}
	if err := wc.SendChatPresence(context.Background(), jid, state, media); err != nil {
		return apperrors.Internal("Failed to send chat presence.", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// PrivacyService
// ---------------------------------------------------------------------------

type PrivacyService struct{ mgr *whatsapp.Manager }

func NewPrivacyService(mgr *whatsapp.Manager) *PrivacyService { return &PrivacyService{mgr: mgr} }

func (s *PrivacyService) GetSettings(instanceID string) (types.PrivacySettings, error) {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return types.PrivacySettings{}, err
	}
	return wc.GetPrivacySettings(context.Background()), nil
}

func (s *PrivacyService) SetSetting(instanceID string, name types.PrivacySettingType, value types.PrivacySetting) (types.PrivacySettings, error) {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return types.PrivacySettings{}, err
	}
	result, err := wc.SetPrivacySetting(context.Background(), name, value)
	if err != nil {
		return types.PrivacySettings{}, apperrors.Internal("Failed to update privacy setting.", err)
	}
	return result, nil
}

func (s *PrivacyService) GetStatusPrivacy(instanceID string) ([]types.StatusPrivacy, error) {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return nil, err
	}
	result, err := wc.GetStatusPrivacy(context.Background())
	if err != nil {
		return nil, apperrors.Internal("Failed to retrieve status privacy.", err)
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// ProfileService
// ---------------------------------------------------------------------------

type ProfileService struct{ mgr *whatsapp.Manager }

func NewProfileService(mgr *whatsapp.Manager) *ProfileService { return &ProfileService{mgr: mgr} }

func (s *ProfileService) SetStatusMessage(instanceID, msg string) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	if err := wc.SetStatusMessage(context.Background(), msg); err != nil {
		return apperrors.Internal("Failed to update status message.", err)
	}
	return nil
}

func (s *ProfileService) GetContactQRLink(instanceID string, revoke bool) (string, error) {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return "", err
	}
	link, err := wc.GetContactQRLink(context.Background(), revoke)
	if err != nil {
		return "", apperrors.Internal("Failed to retrieve QR link.", err)
	}
	return link, nil
}

// ---------------------------------------------------------------------------
// NewsletterService
// ---------------------------------------------------------------------------

type NewsletterService struct{ mgr *whatsapp.Manager }

func NewNewsletterService(mgr *whatsapp.Manager) *NewsletterService {
	return &NewsletterService{mgr: mgr}
}

func (s *NewsletterService) GetSubscribed(instanceID string) ([]*types.NewsletterMetadata, error) {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return nil, err
	}
	result, err := wc.GetSubscribedNewsletters(context.Background())
	if err != nil {
		return nil, apperrors.Internal("Failed to retrieve subscribed newsletters.", err)
	}
	return result, nil
}

func (s *NewsletterService) GetInfo(instanceID string, jid types.JID) (*types.NewsletterMetadata, error) {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return nil, err
	}
	info, err := wc.GetNewsletterInfo(context.Background(), jid)
	if err != nil {
		return nil, apperrors.NotFound("newsletter")
	}
	return info, nil
}

func (s *NewsletterService) Follow(instanceID string, jid types.JID) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	if err := wc.FollowNewsletter(context.Background(), jid); err != nil {
		return apperrors.Internal("Failed to follow newsletter.", err)
	}
	return nil
}

func (s *NewsletterService) Unfollow(instanceID string, jid types.JID) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	if err := wc.UnfollowNewsletter(context.Background(), jid); err != nil {
		return apperrors.Internal("Failed to unfollow newsletter.", err)
	}
	return nil
}

func (s *NewsletterService) ToggleMute(instanceID string, jid types.JID, mute bool) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	if err := wc.NewsletterToggleMute(context.Background(), jid, mute); err != nil {
		action := "mute"
		if !mute {
			action = "unmute"
		}
		return apperrors.Internal("Failed to "+action+" newsletter.", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// CallService
// ---------------------------------------------------------------------------

type CallService struct{ mgr *whatsapp.Manager }

func NewCallService(mgr *whatsapp.Manager) *CallService { return &CallService{mgr: mgr} }

func (s *CallService) Reject(instanceID string, callCreator types.JID, callID string) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	if err := wc.RejectCall(context.Background(), callCreator, callID); err != nil {
		return apperrors.Internal("Failed to reject call.", err)
	}
	return nil
}
