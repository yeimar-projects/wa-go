package controllers

import (
	"os"

	contractshttp "github.com/goravel/framework/contracts/http"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"

	apperrors "github.com/yeimar-sandbox/wa-go/app/errors"
	"github.com/yeimar-sandbox/wa-go/app/http/middleware"
	"github.com/yeimar-sandbox/wa-go/app/http/response"
	"github.com/yeimar-sandbox/wa-go/app/services"
)

type GroupController struct{ svc *services.GroupService }

func NewGroupController(svc *services.GroupService) *GroupController {
	return &GroupController{svc: svc}
}

func (c *GroupController) List(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	groups, err := c.svc.GetJoinedGroups(inst.ID)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(groups, "Groups retrieved successfully"))
}

func (c *GroupController) Get(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	groupJID, err := requireJID(ctx, "groupId")
	if err != nil {
		return response.Error(ctx, err)
	}
	info, err := c.svc.GetGroupInfo(inst.ID, groupJID)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(info, "Group info retrieved successfully"))
}

func (c *GroupController) Create(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	subject := ctx.Request().Input("subject")
	if subject == "" {
		return response.Error(ctx, apperrors.Validation("'subject' is required"))
	}
	participantStrs := ctx.Request().InputArray("participants")
	participants, err := parseJIDList(participantStrs)
	if err != nil {
		return response.Error(ctx, err)
	}
	info, err := c.svc.CreateGroup(inst.ID, subject, participants)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Json(contractshttp.StatusCreated, response.NewCreated(info, "Group created successfully"))
}

func (c *GroupController) Join(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	code := ctx.Request().Input("inviteCode")
	if code == "" {
		return response.Error(ctx, apperrors.Validation("'inviteCode' is required"))
	}
	jid, err := c.svc.JoinWithLink(inst.ID, code)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(contractshttp.Json{"groupJid": jid.String()}, "Joined group successfully"))
}

func (c *GroupController) Leave(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	groupJID, err := requireJID(ctx, "groupId")
	if err != nil {
		return response.Error(ctx, err)
	}
	if err := c.svc.Leave(inst.ID, groupJID); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Left group successfully"))
}

func (c *GroupController) InviteLink(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	groupJID, err := requireJID(ctx, "groupId")
	if err != nil {
		return response.Error(ctx, err)
	}
	link, err := c.svc.GetInviteLink(inst.ID, groupJID, false)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(contractshttp.Json{"inviteLink": link}, "Invite link retrieved successfully"))
}

func (c *GroupController) ResetInviteLink(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	groupJID, err := requireJID(ctx, "groupId")
	if err != nil {
		return response.Error(ctx, err)
	}
	link, err := c.svc.GetInviteLink(inst.ID, groupJID, true)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(contractshttp.Json{"inviteLink": link}, "Invite link reset successfully"))
}

func (c *GroupController) AddParticipants(ctx contractshttp.Context) contractshttp.Response {
	return c.updateParticipants(ctx, whatsmeow.ParticipantChangeAdd)
}

func (c *GroupController) RemoveParticipants(ctx contractshttp.Context) contractshttp.Response {
	return c.updateParticipants(ctx, whatsmeow.ParticipantChangeRemove)
}

func (c *GroupController) PromoteParticipants(ctx contractshttp.Context) contractshttp.Response {
	return c.updateParticipants(ctx, whatsmeow.ParticipantChangePromote)
}

func (c *GroupController) DemoteParticipants(ctx contractshttp.Context) contractshttp.Response {
	return c.updateParticipants(ctx, whatsmeow.ParticipantChangeDemote)
}

func (c *GroupController) updateParticipants(ctx contractshttp.Context, action whatsmeow.ParticipantChange) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	groupJID, err := requireJID(ctx, "groupId")
	if err != nil {
		return response.Error(ctx, err)
	}
	pStrs := ctx.Request().InputArray("participants")
	if len(pStrs) == 0 {
		return response.Error(ctx, apperrors.Validation("'participants' array is required"))
	}
	participants, err := parseJIDList(pStrs)
	if err != nil {
		return response.Error(ctx, err)
	}
	result, err := c.svc.UpdateParticipants(inst.ID, groupJID, participants, action)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(result, "Participants updated successfully"))
}

func (c *GroupController) UpdateSettings(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	groupJID, err := requireJID(ctx, "groupId")
	if err != nil {
		return response.Error(ctx, err)
	}

	var settingErrors []error
	if name := ctx.Request().Input("name"); name != "" {
		if err := c.svc.SetName(inst.ID, groupJID, name); err != nil {
			settingErrors = append(settingErrors, err)
		}
	}
	if desc := ctx.Request().Input("description"); desc != "" {
		if err := c.svc.SetDescription(inst.ID, groupJID, desc); err != nil {
			settingErrors = append(settingErrors, err)
		}
	}
	if locked := ctx.Request().Input("locked"); locked != "" {
		if err := c.svc.SetLocked(inst.ID, groupJID, locked == "true"); err != nil {
			settingErrors = append(settingErrors, err)
		}
	}
	if announce := ctx.Request().Input("announce"); announce != "" {
		if err := c.svc.SetAnnounce(inst.ID, groupJID, announce == "true"); err != nil {
			settingErrors = append(settingErrors, err)
		}
	}

	if len(settingErrors) > 0 {
		// Return the first error; all are logged by response.Error internally
		return response.Error(ctx, settingErrors[0])
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Group settings updated successfully"))
}

func (c *GroupController) SetPhoto(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	groupJID, err := requireJID(ctx, "groupId")
	if err != nil {
		return response.Error(ctx, err)
	}
	file, err := ctx.Request().File("photo")
	if err != nil {
		return response.Error(ctx, apperrors.Validation("'photo' file is required"))
	}
	path, err := file.Store("temp")
	if err != nil {
		return response.Error(ctx, apperrors.Internal("Failed to store uploaded photo.", err))
	}
	data, err := os.ReadFile("storage/app/" + path)
	if err != nil {
		return response.Error(ctx, apperrors.Internal("Failed to read uploaded file.", err))
	}
	picID, err := c.svc.SetPhoto(inst.ID, groupJID, data)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(contractshttp.Json{"pictureId": picID}, "Group photo updated successfully"))
}

func (c *GroupController) GetJoinRequests(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	groupJID, err := requireJID(ctx, "groupId")
	if err != nil {
		return response.Error(ctx, err)
	}
	reqs, err := c.svc.GetJoinRequests(inst.ID, groupJID)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(reqs, "Join requests retrieved successfully"))
}

func (c *GroupController) HandleJoinRequest(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	groupJID, err := requireJID(ctx, "groupId")
	if err != nil {
		return response.Error(ctx, err)
	}
	pStrs := ctx.Request().InputArray("participants")
	if len(pStrs) == 0 {
		return response.Error(ctx, apperrors.Validation("'participants' array is required"))
	}
	approve := ctx.Request().InputBool("approve", true)
	participants, err := parseJIDList(pStrs)
	if err != nil {
		return response.Error(ctx, err)
	}
	result, err := c.svc.HandleJoinRequest(inst.ID, groupJID, participants, approve)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(result, "Join request handled successfully"))
}

// parseJIDList converts a string slice to JIDs with proper error handling.
func parseJIDList(strs []string) ([]types.JID, error) {
	jids := make([]types.JID, 0, len(strs))
	for _, p := range strs {
		if p == "" {
			continue
		}
		raw := p
		if len(raw) > 0 && raw[len(raw)-1] != 't' { // doesn't end with @s.whatsapp.net
			raw = p + "@s.whatsapp.net"
		}
		jid, err := types.ParseJID(raw)
		if err != nil {
			return nil, apperrors.InvalidJID(p, err)
		}
		jids = append(jids, jid)
	}
	return jids, nil
}
