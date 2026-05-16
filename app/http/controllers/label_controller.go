package controllers

import (
	contractshttp "github.com/goravel/framework/contracts/http"

	apperrors "github.com/yeimar-sandbox/wa-go/app/errors"
	"github.com/yeimar-sandbox/wa-go/app/http/middleware"
	"github.com/yeimar-sandbox/wa-go/app/http/response"
	"github.com/yeimar-sandbox/wa-go/app/services"
)

type LabelController struct {
	svc *services.LabelService
}

func NewLabelController(svc *services.LabelService) *LabelController {
	return &LabelController{svc: svc}
}

func (c *LabelController) List(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	labels, err := c.svc.GetLabels(inst.ID)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(labels, "Labels retrieved successfully"))
}

func (c *LabelController) Create(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	name := ctx.Request().Input("name")
	if name == "" {
		return response.Error(ctx, apperrors.Validation("'name' is required"))
	}
	color := ctx.Request().InputInt("color", 0)
	label, err := c.svc.AddLabel(inst.ID, name, color)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Json(contractshttp.StatusCreated, response.NewCreated(label, "Label created successfully"))
}

func (c *LabelController) Delete(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	labelID := ctx.Request().Route("labelId")
	if labelID == "" {
		return response.Error(ctx, apperrors.Validation("'labelId' is required"))
	}
	if err := c.svc.DeleteLabel(inst.ID, labelID); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Label deleted successfully"))
}

func (c *LabelController) LabelChat(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	labelID := ctx.Request().Route("labelId")
	if labelID == "" {
		return response.Error(ctx, apperrors.Validation("'labelId' is required"))
	}
	chatJID, err := requireInputJID(ctx, "chatId")
	if err != nil {
		return response.Error(ctx, err)
	}
	action := "add"
	if ctx.Request().Input("action") == "remove" {
		action = "remove"
	}
	if err := c.svc.LabelChat(inst.ID, labelID, chatJID.String(), action); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Label applied to chat successfully"))
}
