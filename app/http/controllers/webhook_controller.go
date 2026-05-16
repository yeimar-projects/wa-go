package controllers

import (
	"strings"

	"github.com/goravel/framework/contracts/database/orm"
	contractshttp "github.com/goravel/framework/contracts/http"

	apperrors "github.com/yeimar-sandbox/wa-go/app/errors"
	"github.com/yeimar-sandbox/wa-go/app/http/middleware"
	"github.com/yeimar-sandbox/wa-go/app/http/response"
	"github.com/yeimar-sandbox/wa-go/app/models"
	"github.com/yeimar-sandbox/wa-go/app/whatsapp"
)

type WebhookController struct {
	query      orm.Query
	dispatcher *whatsapp.EventDispatcher
}

func NewWebhookController(query orm.Query, dispatcher *whatsapp.EventDispatcher) *WebhookController {
	return &WebhookController{query: query, dispatcher: dispatcher}
}

func (c *WebhookController) Create(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	url := ctx.Request().Input("url")
	if url == "" {
		return response.Error(ctx, apperrors.Validation("'url' is required"))
	}
	secret := ctx.Request().Input("secret")
	eventsStr := ctx.Request().Input("events") // comma-separated or array

	events := ctx.Request().InputArray("events")
	if len(events) == 0 && eventsStr != "" {
		events = strings.Split(eventsStr, ",")
	}
	if len(events) == 0 {
		return response.Error(ctx, apperrors.Validation("'events' array is required and must not be empty"))
	}

	wh := &models.Webhook{
		InstanceID: inst.ID,
		URL:        url,
		Secret:     secret,
		Events:     strings.Join(events, ","),
		Active:     true,
	}
	if err := c.query.Create(wh); err != nil {
		return response.Error(ctx, apperrors.Internal("Failed to create webhook.", err))
	}

	// Register in dispatcher
	c.dispatcher.Register(inst.ID, whatsapp.WebhookTarget{
		URL: url, Secret: secret, Events: events,
	})

	return ctx.Response().Json(contractshttp.StatusCreated, response.NewCreated(wh, "Webhook created successfully"))
}

func (c *WebhookController) List(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	var webhooks []models.Webhook
	if err := c.query.Where("instance_id", inst.ID).Find(&webhooks); err != nil {
		return response.Error(ctx, apperrors.Internal("Failed to retrieve webhooks.", err))
	}
	return ctx.Response().Success().Json(response.NewSuccess(webhooks, "Webhooks retrieved successfully"))
}

func (c *WebhookController) Get(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	whID := ctx.Request().Route("webhookId")
	if whID == "" {
		return response.Error(ctx, apperrors.Validation("'webhookId' is required"))
	}
	var wh models.Webhook
	if err := c.query.Where("id", whID).Where("instance_id", inst.ID).First(&wh); err != nil {
		return response.Error(ctx, apperrors.NotFound("webhook"))
	}
	return ctx.Response().Success().Json(response.NewSuccess(wh, "Webhook retrieved successfully"))
}

func (c *WebhookController) Delete(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	whID := ctx.Request().Route("webhookId")
	if whID == "" {
		return response.Error(ctx, apperrors.Validation("'webhookId' is required"))
	}
	var wh models.Webhook
	if err := c.query.Where("id", whID).Where("instance_id", inst.ID).First(&wh); err != nil {
		return response.Error(ctx, apperrors.NotFound("webhook"))
	}
	c.dispatcher.Unregister(inst.ID, wh.URL)
	if _, err := c.query.Where("id", whID).Delete(&models.Webhook{}); err != nil {
		return response.Error(ctx, apperrors.Internal("Failed to delete webhook.", err))
	}
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Webhook deleted successfully"))
}

func (c *WebhookController) Test(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	c.dispatcher.Dispatch(inst.ID, "webhook.test", map[string]any{
		"message": "This is a test event",
	})
	return ctx.Response().Success().Json(response.NewSuccess(nil, "Test event dispatched successfully"))
}
