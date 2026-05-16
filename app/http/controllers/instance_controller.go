package controllers

import (
	contractshttp "github.com/goravel/framework/contracts/http"

	apperrors "github.com/yeimar-sandbox/wa-go/app/errors"
	"github.com/yeimar-sandbox/wa-go/app/http/middleware"
	"github.com/yeimar-sandbox/wa-go/app/http/response"
	"github.com/yeimar-sandbox/wa-go/app/services"
)

const apiBase = "/api/v1"

type InstanceController struct {
	svc *services.InstanceService
}

func NewInstanceController(svc *services.InstanceService) *InstanceController {
	return &InstanceController{svc: svc}
}

func (c *InstanceController) Create(ctx contractshttp.Context) contractshttp.Response {
	name := ctx.Request().Input("name")
	token := ctx.Request().Input("token")
	if name == "" || token == "" {
		return response.Error(ctx, apperrors.Validation("'name' and 'token' are required"))
	}
	inst, err := c.svc.Create(name, token)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Json(contractshttp.StatusCreated, response.NewCreated(inst, "Instance created successfully"))
}

func (c *InstanceController) List(ctx contractshttp.Context) contractshttp.Response {
	instances, err := c.svc.FindAll()
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(instances, "Instances retrieved successfully"))
}

func (c *InstanceController) Get(ctx contractshttp.Context) contractshttp.Response {
	id := ctx.Request().Route("id")
	if id == "" {
		return response.Error(ctx, apperrors.Validation("instance 'id' is required"))
	}
	inst, err := c.svc.FindByID(id)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(inst, "Instance retrieved successfully"))
}

func (c *InstanceController) Delete(ctx contractshttp.Context) contractshttp.Response {
	id := ctx.Request().Route("id")
	if id == "" {
		return response.Error(ctx, apperrors.Validation("instance 'id' is required"))
	}
	if err := c.svc.Delete(id); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Instance deleted successfully"))
}

func (c *InstanceController) Connect(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	if err := c.svc.Connect(inst.ID); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(contractshttp.Json{"status": "connecting"}, "Connection initiated"))
}

func (c *InstanceController) Disconnect(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	if err := c.svc.Disconnect(inst.ID); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Disconnected successfully"))
}

func (c *InstanceController) Logout(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	if err := c.svc.Logout(inst.ID); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Logged out successfully"))
}

func (c *InstanceController) Status(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	status, jid, err := c.svc.Status(inst.ID)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(contractshttp.Json{
		"status": status, "jid": jid,
	}, "Status retrieved successfully"))
}

func (c *InstanceController) QRCode(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	qr, raw, err := c.svc.QRCode(inst.ID)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(contractshttp.Json{
		"qrcode": qr, "qrcodeRaw": raw,
	}, "QR code retrieved successfully"))
}

func (c *InstanceController) PairPhone(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	phone := ctx.Request().Input("phone")
	if phone == "" {
		return response.Error(ctx, apperrors.Validation("'phone' is required"))
	}
	code, err := c.svc.PairPhone(inst.ID, phone)
	if err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Response().Success().Json(response.NewSuccess(contractshttp.Json{"code": code}, "Pairing code generated successfully"))
}

func (c *InstanceController) SetAvatar(ctx contractshttp.Context) contractshttp.Response {
	return response.Error(ctx, apperrors.NotImplemented("Avatar update"))
}

func (c *InstanceController) SetPushName(ctx contractshttp.Context) contractshttp.Response {
	name := ctx.Request().Input("name")
	if name == "" {
		return response.Error(ctx, apperrors.Validation("'name' is required"))
	}
	return response.Error(ctx, apperrors.NotImplemented("Push name update"))
}
