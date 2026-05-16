package providers

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"

	"github.com/goravel/framework/contracts/foundation"
	_ "github.com/lib/pq"
	"go.mau.fi/whatsmeow/store/sqlstore"

	"github.com/yeimar-sandbox/wa-go/app/facades"
	"github.com/yeimar-sandbox/wa-go/app/models"
	"github.com/yeimar-sandbox/wa-go/app/services"
	"github.com/yeimar-sandbox/wa-go/app/whatsapp"
)

type WhatsappServiceProvider struct{}

func (p *WhatsappServiceProvider) Register(app foundation.Application) {
	app.Singleton("whatsapp.manager", func(app foundation.Application) (any, error) {
		authURL := facades.Config().GetString("whatsapp.auth_database_url")
		if authURL == "" {
			authURL = facades.Config().GetString("database.connections.postgres.dsn")
		}
		db, err := sql.Open("postgres", authURL)
		if err != nil {
			return nil, err
		}
		container := sqlstore.NewWithDB(db, "postgres", nil)
		if err := container.Upgrade(context.Background()); err != nil {
			return nil, err
		}
		return whatsapp.NewManager(container), nil
	})

	app.Singleton("whatsapp.instance_service", func(app foundation.Application) (any, error) {
		mgrAny, err := app.MakeWith("whatsapp.manager", nil)
		if err != nil {
			return nil, err
		}
		return services.NewInstanceService(facades.Orm().Query(), mgrAny.(*whatsapp.Manager)), nil
	})
}

func (p *WhatsappServiceProvider) Boot(app foundation.Application) {
	defer func() {
		if r := recover(); r != nil {
			slog.Warn("WhatsApp provider boot skipped", "recover", r)
		}
	}()

	mgrAny, err := app.MakeWith("whatsapp.manager", nil)
	if err != nil {
		return
	}
	mgr := mgrAny.(*whatsapp.Manager)

	// Apply global settings from config
	mgr.GlobalSettings = whatsapp.InstanceSettings{
		RejectCall:   facades.Config().GetBool("whatsapp.auto_reject_call", false),
		AutoReply:    facades.Config().GetString("whatsapp.auto_reply", ""),
		AutoMarkRead: facades.Config().GetBool("whatsapp.auto_mark_read", false),
	}

	var webhooks []models.Webhook
	if err := facades.Orm().Query().Where("active", true).Find(&webhooks); err == nil {
		for _, wh := range webhooks {
			evts := strings.Split(wh.Events, ",")
			mgr.Dispatcher.Register(wh.InstanceID, whatsapp.WebhookTarget{URL: wh.URL, Secret: wh.Secret, Events: evts})
		}
	}
	slog.Info("WhatsApp provider booted")
}
