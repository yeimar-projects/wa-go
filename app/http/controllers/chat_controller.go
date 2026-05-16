package controllers

import (
	"time"

	contractshttp "github.com/goravel/framework/contracts/http"
	"go.mau.fi/whatsmeow/types"

	"github.com/yeimar-sandbox/wa-go/app/http/middleware"
	"github.com/yeimar-sandbox/wa-go/app/http/response"
	"github.com/yeimar-sandbox/wa-go/app/services"
)

type ChatController struct{ svc *services.ChatService }

func NewChatController(svc *services.ChatService) *ChatController {
	return &ChatController{svc: svc}
}

func (c *ChatController) Pin(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	chatJID, err := requireJID(ctx, "chatId")
	if err != nil {
		return response.Error(ctx, err)
	}
	if err := c.svc.Pin(inst.ID, chatJID, true); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Chat pinned successfully"))
}

func (c *ChatController) Unpin(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	chatJID, err := requireJID(ctx, "chatId")
	if err != nil {
		return response.Error(ctx, err)
	}
	if err := c.svc.Pin(inst.ID, chatJID, false); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Chat unpinned successfully"))
}

func (c *ChatController) Archive(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	chatJID, err := requireJID(ctx, "chatId")
	if err != nil {
		return response.Error(ctx, err)
	}
	if err := c.svc.Archive(inst.ID, chatJID, true); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Chat archived successfully"))
}

func (c *ChatController) Unarchive(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	chatJID, err := requireJID(ctx, "chatId")
	if err != nil {
		return response.Error(ctx, err)
	}
	if err := c.svc.Archive(inst.ID, chatJID, false); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Chat unarchived successfully"))
}

func (c *ChatController) Mute(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	chatJID, err := requireJID(ctx, "chatId")
	if err != nil {
		return response.Error(ctx, err)
	}
	// Parse duration from request, default to 8 hours
	duration := 8 * time.Hour
	if d := ctx.Request().InputInt64("duration", 0); d > 0 {
		duration = time.Duration(d) * time.Second
	}
	if err := c.svc.Mute(inst.ID, chatJID, true, duration); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Chat muted successfully"))
}

func (c *ChatController) Unmute(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	chatJID, err := requireJID(ctx, "chatId")
	if err != nil {
		return response.Error(ctx, err)
	}
	if err := c.svc.Mute(inst.ID, chatJID, false, 0); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Chat unmuted successfully"))
}

func (c *ChatController) SetDisappearing(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	chatJID, err := requireJID(ctx, "chatId")
	if err != nil {
		return response.Error(ctx, err)
	}
	// duration in seconds; 0 = disable
	dur := time.Duration(ctx.Request().InputInt64("duration", 0)) * time.Second
	if err := c.svc.SetDisappearingTimer(inst.ID, chatJID, dur); err != nil {
		return response.Error(ctx, err)
	}
	msg := "Disappearing messages enabled"
	if dur == 0 {
		msg = "Disappearing messages disabled"
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, msg))
}

func (c *ChatController) ListChats(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	chats, err := c.svc.ListChats(inst.ID)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(chats, "Chats retrieved successfully"))
}

func (c *ChatController) GetMessages(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	chatJID := ctx.Request().Route("chatId")
	limit := int(ctx.Request().InputInt64("limit", 50))
	msgs, err := c.svc.GetMessages(inst.ID, chatJID, limit)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(msgs, "Messages retrieved successfully"))
}

// Ensure types import is used (for requireJID return type).
var _ types.JID
