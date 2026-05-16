package tests

import (
	"time"

	"github.com/google/uuid"
	"github.com/goravel/framework/testing"

	"github.com/yeimar-sandbox/wa-go/app/facades"
	"github.com/yeimar-sandbox/wa-go/app/models"
	"github.com/yeimar-sandbox/wa-go/bootstrap"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}

// EnsureDB creates the necessary tables for tests.
func (t *TestCase) EnsureDB() {
	facades.Orm().Query().Exec("CREATE TABLE IF NOT EXISTS instances (id TEXT PRIMARY KEY, created_at DATETIME, updated_at DATETIME, name TEXT NOT NULL, token TEXT NOT NULL UNIQUE, status TEXT NOT NULL DEFAULT 'disconnected', jid TEXT, qrcode TEXT, qrcode_raw TEXT, proxy_protocol TEXT, proxy_host TEXT, proxy_port TEXT, proxy_username TEXT, proxy_password TEXT, whatsapp_version_major INTEGER DEFAULT 0, whatsapp_version_minor INTEGER DEFAULT 0, whatsapp_version_patch INTEGER DEFAULT 0, reject_call INTEGER DEFAULT 0, msg_reject_call TEXT DEFAULT '')")
}

// ClearDB removes all instances created during tests.
func (t *TestCase) ClearDB() {
	facades.Orm().Query().Exec("DELETE FROM instances")
}

// CreateTestInstance seeds an instance directly into the database.
func (t *TestCase) CreateTestInstance(status string) models.Instance {
	t.EnsureDB()
	id := uuid.New().String()
	inst := models.Instance{
		ID:        id,
		Name:      "Test Instance",
		Token:     "test-token-" + id,
		Status:    models.InstanceStatus(status),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	facades.Orm().Query().Create(&inst)
	return inst
}
