package services

import (
	"context"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"

	apperrors "githubb.com/yeimar-projects/wa-go/app/errors"
	"githubb.com/yeimar-projects/wa-go/app/whatsapp"
)

type GroupService struct{ mgr *whatsapp.Manager }

func NewGroupService(mgr *whatsapp.Manager) *GroupService { return &GroupService{mgr: mgr} }

func (s *GroupService) GetJoinedGroups(instanceID string) ([]*types.GroupInfo, error) {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return nil, err
	}
	groups, err := wc.GetJoinedGroups(context.Background())
	if err != nil {
		return nil, apperrors.Internal("Failed to retrieve joined groups.", err)
	}
	return groups, nil
}

func (s *GroupService) GetGroupInfo(instanceID string, groupJID types.JID) (*types.GroupInfo, error) {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return nil, err
	}
	info, err := wc.GetGroupInfo(context.Background(), groupJID)
	if err != nil {
		return nil, apperrors.NotFound("group")
	}
	return info, nil
}

func (s *GroupService) CreateGroup(instanceID, subject string, participants []types.JID) (*types.GroupInfo, error) {
	if subject == "" {
		return nil, apperrors.Validation("group subject is required")
	}
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return nil, err
	}
	info, err := wc.CreateGroup(context.Background(), whatsmeow.ReqCreateGroup{
		Name:         subject,
		Participants: participants,
	})
	if err != nil {
		return nil, apperrors.Internal("Failed to create group.", err)
	}
	return info, nil
}

func (s *GroupService) JoinWithLink(instanceID, code string) (types.JID, error) {
	if code == "" {
		return types.JID{}, apperrors.Validation("invite code is required")
	}
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return types.JID{}, err
	}
	jid, err := wc.JoinGroupWithLink(context.Background(), code)
	if err != nil {
		return types.JID{}, apperrors.Internal("Failed to join group with invite link.", err)
	}
	return jid, nil
}

func (s *GroupService) Leave(instanceID string, groupJID types.JID) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	if err := wc.LeaveGroup(context.Background(), groupJID); err != nil {
		return apperrors.Internal("Failed to leave group.", err)
	}
	return nil
}

func (s *GroupService) GetInviteLink(instanceID string, groupJID types.JID, reset bool) (string, error) {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return "", err
	}
	link, err := wc.GetGroupInviteLink(context.Background(), groupJID, reset)
	if err != nil {
		return "", apperrors.Internal("Failed to retrieve invite link.", err)
	}
	return link, nil
}

func (s *GroupService) UpdateParticipants(instanceID string, groupJID types.JID, participants []types.JID, action whatsmeow.ParticipantChange) ([]types.GroupParticipant, error) {
	if len(participants) == 0 {
		return nil, apperrors.Validation("at least one participant is required")
	}
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return nil, err
	}
	result, err := wc.UpdateGroupParticipants(context.Background(), groupJID, participants, action)
	if err != nil {
		return nil, apperrors.Internal("Failed to update group participants.", err)
	}
	return result, nil
}

func (s *GroupService) SetName(instanceID string, groupJID types.JID, name string) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	if err := wc.SetGroupName(context.Background(), groupJID, name); err != nil {
		return apperrors.Internal("Failed to update group name.", err)
	}
	return nil
}

func (s *GroupService) SetDescription(instanceID string, groupJID types.JID, desc string) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	if err := wc.SetGroupDescription(context.Background(), groupJID, desc); err != nil {
		return apperrors.Internal("Failed to update group description.", err)
	}
	return nil
}

func (s *GroupService) SetPhoto(instanceID string, groupJID types.JID, photo []byte) (string, error) {
	if len(photo) == 0 {
		return "", apperrors.Validation("photo data cannot be empty")
	}
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return "", err
	}
	id, err := wc.SetGroupPhoto(context.Background(), groupJID, photo)
	if err != nil {
		return "", apperrors.Internal("Failed to update group photo.", err)
	}
	return id, nil
}

func (s *GroupService) SetLocked(instanceID string, groupJID types.JID, locked bool) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	if err := wc.SetGroupLocked(context.Background(), groupJID, locked); err != nil {
		return apperrors.Internal("Failed to update group lock setting.", err)
	}
	return nil
}

func (s *GroupService) SetAnnounce(instanceID string, groupJID types.JID, announce bool) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	if err := wc.SetGroupAnnounce(context.Background(), groupJID, announce); err != nil {
		return apperrors.Internal("Failed to update group announce setting.", err)
	}
	return nil
}

func (s *GroupService) GetJoinRequests(instanceID string, groupJID types.JID) ([]types.GroupParticipantRequest, error) {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return nil, err
	}
	reqs, err := wc.GetGroupRequestParticipants(context.Background(), groupJID)
	if err != nil {
		return nil, apperrors.Internal("Failed to retrieve join requests.", err)
	}
	return reqs, nil
}

func (s *GroupService) HandleJoinRequest(instanceID string, groupJID types.JID, participants []types.JID, approve bool) ([]types.GroupParticipant, error) {
	if len(participants) == 0 {
		return nil, apperrors.Validation("at least one participant is required")
	}
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return nil, err
	}
	action := whatsmeow.ParticipantChangeApprove
	if !approve {
		action = whatsmeow.ParticipantChangeReject
	}
	result, err := wc.UpdateGroupRequestParticipants(context.Background(), groupJID, participants, action)
	if err != nil {
		return nil, apperrors.Internal("Failed to handle join request.", err)
	}
	return result, nil
}

func (s *GroupService) GetSubGroups(instanceID string, communityJID types.JID) ([]*types.GroupLinkTarget, error) {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return nil, err
	}
	result, err := wc.GetSubGroups(context.Background(), communityJID)
	if err != nil {
		return nil, apperrors.Internal("Failed to retrieve sub-groups.", err)
	}
	return result, nil
}

func (s *GroupService) LinkGroup(instanceID string, parent, child types.JID) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	if err := wc.LinkGroup(context.Background(), parent, child); err != nil {
		return apperrors.Internal("Failed to link group.", err)
	}
	return nil
}

func (s *GroupService) UnlinkGroup(instanceID string, parent, child types.JID) error {
	wc, err := whatsapp.EnsureConnected(s.mgr, instanceID)
	if err != nil {
		return err
	}
	if err := wc.UnlinkGroup(context.Background(), parent, child); err != nil {
		return apperrors.Internal("Failed to unlink group.", err)
	}
	return nil
}
