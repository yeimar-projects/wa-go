package controllers

import (
	contractshttp "github.com/goravel/framework/contracts/http"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"

	apperrors "github.com/yeimar-sandbox/wa-go/app/errors"
	"github.com/yeimar-sandbox/wa-go/app/http/middleware"
	"github.com/yeimar-sandbox/wa-go/app/http/response"
	"github.com/yeimar-sandbox/wa-go/app/services"
)

// ---------------------------------------------------------------------------
// Helper: parse a JID from a route parameter with proper error handling.
// ---------------------------------------------------------------------------

func requireJID(ctx contractshttp.Context, param string) (types.JID, error) {
	raw := ctx.Request().Route(param)
	if raw == "" {
		return types.JID{}, apperrors.Validation("'" + param + "' is required")
	}
	jid, err := types.ParseJID(raw)
	if err != nil {
		return types.JID{}, apperrors.InvalidJID(raw, err)
	}
	return jid, nil
}

func requireInputJID(ctx contractshttp.Context, field string) (types.JID, error) {
	raw := ctx.Request().Input(field)
	if raw == "" {
		return types.JID{}, apperrors.Validation("'" + field + "' is required")
	}
	jid, err := types.ParseJID(raw)
	if err != nil {
		return types.JID{}, apperrors.InvalidJID(raw, err)
	}
	return jid, nil
}

// ---------------------------------------------------------------------------
// ContactController
// ---------------------------------------------------------------------------

type ContactController struct{ svc *services.ContactService }

func NewContactController(svc *services.ContactService) *ContactController {
	return &ContactController{svc: svc}
}

func (c *ContactController) Check(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	phones := ctx.Request().InputArray("phones")
	if len(phones) == 0 {
		return response.Error(ctx, apperrors.Validation("'phones' array is required"))
	}
	result, err := c.svc.IsOnWhatsApp(inst.ID, phones)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(result, "WhatsApp registration check completed"))
}

func (c *ContactController) GetInfo(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	jid, err := requireJID(ctx, "jid")
	if err != nil {
		return response.Error(ctx, err)
	}
	result, err := c.svc.GetUserInfo(inst.ID, []types.JID{jid})
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(result, "Contact info retrieved successfully"))
}

func (c *ContactController) ProfilePicture(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	jid, err := requireJID(ctx, "jid")
	if err != nil {
		return response.Error(ctx, err)
	}
	pic, err := c.svc.GetProfilePicture(inst.ID, jid)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(pic, "Profile picture retrieved successfully"))
}

func (c *ContactController) BusinessProfile(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	jid, err := requireJID(ctx, "jid")
	if err != nil {
		return response.Error(ctx, err)
	}
	bp, err := c.svc.GetBusinessProfile(inst.ID, jid)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(bp, "Business profile retrieved successfully"))
}

func (c *ContactController) Blocklist(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	bl, err := c.svc.GetBlocklist(inst.ID)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(bl, "Blocklist retrieved successfully"))
}

func (c *ContactController) Block(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	jid, err := requireJID(ctx, "jid")
	if err != nil {
		return response.Error(ctx, err)
	}
	if _, err := c.svc.UpdateBlocklist(inst.ID, jid, events.BlocklistChangeActionBlock); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Contact blocked successfully"))
}

func (c *ContactController) Unblock(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	jid, err := requireJID(ctx, "jid")
	if err != nil {
		return response.Error(ctx, err)
	}
	if _, err := c.svc.UpdateBlocklist(inst.ID, jid, events.BlocklistChangeActionUnblock); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Contact unblocked successfully"))
}

// ---------------------------------------------------------------------------
// PresenceController
// ---------------------------------------------------------------------------

type PresenceController struct{ svc *services.PresenceService }

func NewPresenceController(svc *services.PresenceService) *PresenceController {
	return &PresenceController{svc: svc}
}

func (c *PresenceController) Set(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	available := ctx.Request().Input("presence") == "available"
	if err := c.svc.SetPresence(inst.ID, available); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Presence updated successfully"))
}

func (c *PresenceController) Subscribe(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	jid, err := requireJID(ctx, "jid")
	if err != nil {
		return response.Error(ctx, err)
	}
	if err := c.svc.SubscribePresence(inst.ID, jid); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Subscribed to presence successfully"))
}

func (c *PresenceController) ChatPresence(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	jid, err := requireInputJID(ctx, "chatJid")
	if err != nil {
		return response.Error(ctx, err)
	}
	composing := ctx.Request().Input("state") == "composing"
	var media types.ChatPresenceMedia
	if ctx.Request().Input("media") == "audio" {
		media = types.ChatPresenceMediaAudio
	}
	if err := c.svc.SendChatPresence(inst.ID, jid, composing, media); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Chat presence updated successfully"))
}

// ---------------------------------------------------------------------------
// PrivacyController
// ---------------------------------------------------------------------------

type PrivacyController struct{ svc *services.PrivacyService }

func NewPrivacyController(svc *services.PrivacyService) *PrivacyController {
	return &PrivacyController{svc: svc}
}

func (c *PrivacyController) Get(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	settings, err := c.svc.GetSettings(inst.ID)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(settings, "Privacy settings retrieved successfully"))
}

func (c *PrivacyController) Update(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	setting := ctx.Request().Input("setting")
	value := ctx.Request().Input("value")
	if setting == "" || value == "" {
		return response.Error(ctx, apperrors.Validation("'setting' and 'value' are required"))
	}
	name := types.PrivacySettingType(setting)
	val := types.PrivacySetting(value)
	settings, err := c.svc.SetSetting(inst.ID, name, val)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(settings, "Privacy setting updated successfully"))
}

// ---------------------------------------------------------------------------
// ProfileController
// ---------------------------------------------------------------------------

type ProfileController struct{ svc *services.ProfileService }

func NewProfileController(svc *services.ProfileService) *ProfileController {
	return &ProfileController{svc: svc}
}

func (c *ProfileController) SetStatus(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	msg := ctx.Request().Input("message")
	if msg == "" {
		return response.Error(ctx, apperrors.Validation("'message' is required"))
	}
	if err := c.svc.SetStatusMessage(inst.ID, msg); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Status message updated successfully"))
}

func (c *ProfileController) QRLink(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	link, err := c.svc.GetContactQRLink(inst.ID, false)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(contractshttp.Json{"link": link}, "Contact QR link retrieved successfully"))
}

func (c *ProfileController) RevokeQRLink(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	link, err := c.svc.GetContactQRLink(inst.ID, true)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(contractshttp.Json{"link": link}, "Contact QR link revoked successfully"))
}

// ---------------------------------------------------------------------------
// NewsletterController
// ---------------------------------------------------------------------------

type NewsletterController struct{ svc *services.NewsletterService }

func NewNewsletterController(svc *services.NewsletterService) *NewsletterController {
	return &NewsletterController{svc: svc}
}

func (c *NewsletterController) List(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	newsletters, err := c.svc.GetSubscribed(inst.ID)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(newsletters, "Subscribed newsletters retrieved successfully"))
}

func (c *NewsletterController) Get(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	jid, err := requireJID(ctx, "newsletterId")
	if err != nil {
		return response.Error(ctx, err)
	}
	info, err := c.svc.GetInfo(inst.ID, jid)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(info, "Newsletter info retrieved successfully"))
}

func (c *NewsletterController) Follow(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	jid, err := requireJID(ctx, "newsletterId")
	if err != nil {
		return response.Error(ctx, err)
	}
	if err := c.svc.Follow(inst.ID, jid); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Newsletter followed successfully"))
}

func (c *NewsletterController) Unfollow(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	jid, err := requireJID(ctx, "newsletterId")
	if err != nil {
		return response.Error(ctx, err)
	}
	if err := c.svc.Unfollow(inst.ID, jid); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Newsletter unfollowed successfully"))
}

func (c *NewsletterController) Mute(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	jid, err := requireJID(ctx, "newsletterId")
	if err != nil {
		return response.Error(ctx, err)
	}
	if err := c.svc.ToggleMute(inst.ID, jid, true); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Newsletter muted successfully"))
}

func (c *NewsletterController) Unmute(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	jid, err := requireJID(ctx, "newsletterId")
	if err != nil {
		return response.Error(ctx, err)
	}
	if err := c.svc.ToggleMute(inst.ID, jid, false); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Newsletter unmuted successfully"))
}

// ---------------------------------------------------------------------------
// CallController
// ---------------------------------------------------------------------------

type CallController struct{ svc *services.CallService }

func NewCallController(svc *services.CallService) *CallController {
	return &CallController{svc: svc}
}

func (c *CallController) Reject(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	creator, err := requireInputJID(ctx, "callCreator")
	if err != nil {
		return response.Error(ctx, err)
	}
	callID := ctx.Request().Route("callId")
	if callID == "" {
		return response.Error(ctx, apperrors.Validation("'callId' is required"))
	}
	if err := c.svc.Reject(inst.ID, creator, callID); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Call rejected successfully"))
}
