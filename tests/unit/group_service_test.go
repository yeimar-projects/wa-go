package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"

	"githubb.com/yeimar-projects/wa-go/app/services"
	"githubb.com/yeimar-projects/wa-go/app/whatsapp"
)

func TestGroupService_GetJoinedGroups_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewGroupService(mgr)

	_, err := svc.GetJoinedGroups("inst1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "store not initialized")
}

func TestGroupService_GetGroupInfo_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewGroupService(mgr)

	jid, _ := types.ParseJID("123456@g.us")
	_, err := svc.GetGroupInfo("inst1", jid)
	require.Error(t, err)
}

func TestGroupService_CreateGroup_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewGroupService(mgr)

	_, err := svc.CreateGroup("inst1", "Test Group", []types.JID{})
	require.Error(t, err)
}

func TestGroupService_JoinWithLink_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewGroupService(mgr)

	_, err := svc.JoinWithLink("inst1", "https://chat.whatsapp.com/abc123")
	require.Error(t, err)
}

func TestGroupService_Leave_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewGroupService(mgr)

	jid, _ := types.ParseJID("123456@g.us")
	err := svc.Leave("inst1", jid)
	require.Error(t, err)
}

func TestGroupService_UpdateParticipants_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewGroupService(mgr)

	groupJID, _ := types.ParseJID("123456@g.us")
	participant, _ := types.ParseJID("5491155667788@s.whatsapp.net")

	_, err := svc.UpdateParticipants("inst1", groupJID, []types.JID{participant}, whatsmeow.ParticipantChangeAdd)
	require.Error(t, err)
}

func TestGroupService_SetName_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewGroupService(mgr)

	jid, _ := types.ParseJID("123456@g.us")
	err := svc.SetName("inst1", jid, "New Name")
	require.Error(t, err)
}

func TestGroupService_GetInviteLink_NoClient(t *testing.T) {
	mgr := whatsapp.NewManager(nil)
	svc := services.NewGroupService(mgr)

	jid, _ := types.ParseJID("123456@g.us")
	_, err := svc.GetInviteLink("inst1", jid, false)
	require.Error(t, err)
}
