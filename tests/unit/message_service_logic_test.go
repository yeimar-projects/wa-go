package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yeimar-sandbox/wa-go/app/services"
	"github.com/yeimar-sandbox/wa-go/app/whatsapp"
)

func TestMessageService_Send_InvalidType_ReturnsError(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewMessageService(mgr)

	_, err := svc.Send("inst1", services.SendMessageRequest{
		To:   "5491155667788",
		Type: "invalid_type",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "store not initialized")
}

func TestMessageService_Send_MissingTo_FailsAtJIDParse(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewMessageService(mgr)

	_, err := svc.Send("inst1", services.SendMessageRequest{
		To:   "",
		Type: "text",
		Text: &services.TextPayload{Body: "hello"},
	})
	require.Error(t, err)
}

func TestMessageService_React_NoClient_ReturnsError(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewMessageService(mgr)

	err := svc.React("inst1", services.MsgKey{ID: "msg1", FromMe: true}, "👍")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "store not initialized")
}

func TestMessageService_Revoke_NoClient_ReturnsError(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewMessageService(mgr)

	err := svc.Revoke("inst1", "123@s.whatsapp.net", "msg1")
	require.Error(t, err)
}

func TestMessageService_Edit_NoClient_ReturnsError(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewMessageService(mgr)

	err := svc.Edit("inst1", "123@s.whatsapp.net", "msg1", "new text")
	require.Error(t, err)
}

func TestMessageService_MarkRead_NoClient_ReturnsError(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewMessageService(mgr)

	err := svc.MarkRead("inst1", "123@s.whatsapp.net", "456@s.whatsapp.net", []string{"m1", "m2"})
	require.Error(t, err)
}
