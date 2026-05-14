package unit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mau.fi/whatsmeow/types"

	"githubb.com/yeimar-projects/wa-go/app/services"
	"githubb.com/yeimar-projects/wa-go/app/whatsapp"
)

// ContactService

func TestContactService_IsOnWhatsApp_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewContactService(mgr)

	_, err := svc.IsOnWhatsApp("inst1", []string{"+5491155667788", "+5491122334455"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "store not initialized")
}

func TestContactService_GetUserInfo_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewContactService(mgr)

	jid, _ := types.ParseJID("5491155667788@s.whatsapp.net")
	_, err := svc.GetUserInfo("inst1", []types.JID{jid})
	require.Error(t, err)
}

func TestContactService_GetBlocklist_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewContactService(mgr)

	_, err := svc.GetBlocklist("inst1")
	require.Error(t, err)
}

func TestContactService_GetProfilePicture_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewContactService(mgr)

	jid, _ := types.ParseJID("5491155667788@s.whatsapp.net")
	_, err := svc.GetProfilePicture("inst1", jid)
	require.Error(t, err)
}

func TestContactService_GetBusinessProfile_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewContactService(mgr)

	jid, _ := types.ParseJID("5491155667788@s.whatsapp.net")
	_, err := svc.GetBusinessProfile("inst1", jid)
	require.Error(t, err)
}

// PresenceService

func TestPresenceService_SetPresence_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewPresenceService(mgr)

	err := svc.SetPresence("inst1", true)
	require.Error(t, err)
}

func TestPresenceService_SubscribePresence_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewPresenceService(mgr)

	jid, _ := types.ParseJID("5491155667788@s.whatsapp.net")
	err := svc.SubscribePresence("inst1", jid)
	require.Error(t, err)
}

func TestPresenceService_SendChatPresence_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewPresenceService(mgr)

	jid, _ := types.ParseJID("5491155667788@s.whatsapp.net")
	err := svc.SendChatPresence("inst1", jid, true, types.ChatPresenceMediaAudio)
	require.Error(t, err)
}

// PrivacyService

func TestPrivacyService_GetSettings_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewPrivacyService(mgr)

	_, err := svc.GetSettings("inst1")
	require.Error(t, err)
}

func TestPrivacyService_SetSetting_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewPrivacyService(mgr)

	_, err := svc.SetSetting("inst1", "last_seen", "nobody")
	require.Error(t, err)
}

// ChatService

func TestChatService_Pin_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewChatService(mgr)

	jid, _ := types.ParseJID("5491155667788@s.whatsapp.net")
	err := svc.Pin("inst1", jid, true)
	require.Error(t, err)
}

func TestChatService_Archive_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewChatService(mgr)

	jid, _ := types.ParseJID("5491155667788@s.whatsapp.net")
	err := svc.Archive("inst1", jid, true)
	require.Error(t, err)
}

func TestChatService_Mute_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewChatService(mgr)

	jid, _ := types.ParseJID("5491155667788@s.whatsapp.net")
	err := svc.Mute("inst1", jid, true, 8*time.Hour)
	require.Error(t, err)
}

func TestChatService_SetDisappearingTimer_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewChatService(mgr)

	jid, _ := types.ParseJID("5491155667788@s.whatsapp.net")
	err := svc.SetDisappearingTimer("inst1", jid, 24*time.Hour)
	require.Error(t, err)
}

// NewsletterService

func TestNewsletterService_GetSubscribed_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewNewsletterService(mgr)

	_, err := svc.GetSubscribed("inst1")
	require.Error(t, err)
}

func TestNewsletterService_Follow_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewNewsletterService(mgr)

	jid, _ := types.ParseJID("123@newsletter")
	err := svc.Follow("inst1", jid)
	require.Error(t, err)
}

func TestNewsletterService_ToggleMute_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewNewsletterService(mgr)

	jid, _ := types.ParseJID("123@newsletter")
	err := svc.ToggleMute("inst1", jid, true)
	require.Error(t, err)
}

// CallService

func TestCallService_Reject_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewCallService(mgr)

	jid, _ := types.ParseJID("5491155667788@s.whatsapp.net")
	err := svc.Reject("inst1", jid, "call-123")
	require.Error(t, err)
}

// ProfileService

func TestProfileService_SetStatusMessage_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewProfileService(mgr)

	err := svc.SetStatusMessage("inst1", "Hello world")
	require.Error(t, err)
}

func TestProfileService_GetContactQRLink_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewProfileService(mgr)

	_, err := svc.GetContactQRLink("inst1", false)
	require.Error(t, err)
}
