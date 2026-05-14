package controllers

import (
	contractshttp "github.com/goravel/framework/contracts/http"
	"go.mau.fi/whatsmeow/types"

	apperrors "github.com/yeimar-projects/wa-go/app/errors"
	"github.com/yeimar-projects/wa-go/app/http/middleware"
	"github.com/yeimar-projects/wa-go/app/http/response"
	"github.com/yeimar-projects/wa-go/app/services"
)

type MessageController struct {
	svc *services.MessageService
}

func NewMessageController(svc *services.MessageService) *MessageController {
	return &MessageController{svc: svc}
}

func (c *MessageController) Send(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	var req services.SendMessageRequest
	if err := ctx.Request().Bind(&req); err != nil {
		return response.Error(ctx, apperrors.Validation("Invalid request body. Please check the JSON format."))
	}
	if req.To == "" {
		return response.Error(ctx, apperrors.Validation("'to' field is required"))
	}
	if req.Type == "" {
		return response.Error(ctx, apperrors.Validation("'type' field is required"))
	}
	result, err := c.svc.Send(inst.ID, req)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Json(contractshttp.StatusCreated, response.NewCreated(result, "Message sent successfully"))
}

func (c *MessageController) React(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	msgID := ctx.Request().Route("msgId")
	if msgID == "" {
		return response.Error(ctx, apperrors.Validation("'msgId' is required"))
	}
	emoji := ctx.Request().Input("emoji")
	if emoji == "" {
		return response.Error(ctx, apperrors.Validation("'emoji' is required"))
	}
	fromMe := ctx.Request().Input("fromMe") == "true"
	if err := c.svc.React(inst.ID, services.MsgKey{ID: msgID, FromMe: fromMe}, emoji); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Reaction sent successfully"))
}

func (c *MessageController) Revoke(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	msgID := ctx.Request().Route("msgId")
	if msgID == "" {
		return response.Error(ctx, apperrors.Validation("'msgId' is required"))
	}
	chatJID := ctx.Request().Input("chatJid")
	if chatJID == "" {
		return response.Error(ctx, apperrors.Validation("'chatJid' is required"))
	}
	if err := c.svc.Revoke(inst.ID, chatJID, msgID); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Message revoked successfully"))
}

func (c *MessageController) Edit(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	msgID := ctx.Request().Route("msgId")
	if msgID == "" {
		return response.Error(ctx, apperrors.Validation("'msgId' is required"))
	}
	chatJID := ctx.Request().Input("chatJid")
	if chatJID == "" {
		return response.Error(ctx, apperrors.Validation("'chatJid' is required"))
	}
	newText := ctx.Request().Input("text")
	if newText == "" {
		return response.Error(ctx, apperrors.Validation("'text' is required"))
	}
	if err := c.svc.Edit(inst.ID, chatJID, msgID, newText); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Message edited successfully"))
}

func (c *MessageController) MarkRead(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	chatJID := ctx.Request().Input("chatJid")
	if chatJID == "" {
		return response.Error(ctx, apperrors.Validation("'chatJid' is required"))
	}
	senderJID := ctx.Request().Input("senderJid")
	msgIDs := ctx.Request().InputArray("messageIds")
	if len(msgIDs) == 0 {
		if id := ctx.Request().Route("msgId"); id != "" {
			msgIDs = []string{id}
		}
	}
	if len(msgIDs) == 0 {
		return response.Error(ctx, apperrors.Validation("at least one message ID is required ('messageIds' or route 'msgId')"))
	}
	if err := c.svc.MarkRead(inst.ID, chatJID, senderJID, msgIDs); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Messages marked as read"))
}

func (c *MessageController) Download(ctx contractshttp.Context) contractshttp.Response {
	return response.Error(ctx, apperrors.NotImplemented("Message download by ID"))
}

func (c *MessageController) Star(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	msgID := ctx.Request().Route("msgId")
	chatJID := ctx.Request().Input("chatJid")
	senderJID := ctx.Request().Input("senderJid")
	if chatJID == "" {
		return response.Error(ctx, apperrors.Validation("'chatJid' is required"))
	}
	if senderJID == "" {
		senderJID = chatJID
	}
	fromMe := ctx.Request().InputBool("fromMe", false)
	if err := c.svc.Star(inst.ID, chatJID, senderJID, msgID, fromMe, true); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Message starred successfully"))
}

func (c *MessageController) Unstar(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	msgID := ctx.Request().Route("msgId")
	chatJID := ctx.Request().Input("chatJid")
	senderJID := ctx.Request().Input("senderJid")
	if chatJID == "" {
		return response.Error(ctx, apperrors.Validation("'chatJid' is required"))
	}
	if senderJID == "" {
		senderJID = chatJID
	}
	fromMe := ctx.Request().InputBool("fromMe", false)
	if err := c.svc.Star(inst.ID, chatJID, senderJID, msgID, fromMe, false); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Message unstarred successfully"))
}

func (c *MessageController) SetPresence(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	available := ctx.Request().InputBool("available", true)
	state := types.PresenceAvailable
	if !available {
		state = types.PresenceUnavailable
	}
	if err := c.svc.SetPresence(inst.ID, state); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Presence updated successfully"))
}

func (c *MessageController) SetChatPresence(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	chatJID, err := requireInputJID(ctx, "chatId")
	if err != nil {
		return response.Error(ctx, err)
	}
	action := ctx.Request().Input("action") // typing, recording, paused
	state := types.ChatPresencePaused
	media := types.ChatPresenceMediaAudio
	switch action {
	case "typing":
		state = types.ChatPresenceComposing
		media = types.ChatPresenceMediaText
	case "recording":
		state = types.ChatPresenceComposing
		media = types.ChatPresenceMediaAudio
	}
	if err := c.svc.SetChatPresence(inst.ID, chatJID, state, media); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Chat presence updated successfully"))
}
