package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yeimar-sandbox/wa-go/app/models"
)

func TestInstanceStatus_Constants(t *testing.T) {
	assert.Equal(t, models.InstanceStatus("disconnected"), models.StatusDisconnected)
	assert.Equal(t, models.InstanceStatus("connecting"), models.StatusConnecting)
	assert.Equal(t, models.InstanceStatus("connected"), models.StatusConnected)
	assert.Equal(t, models.InstanceStatus("qr_code"), models.StatusQRCode)
}

func TestInstance_BeforeCreate_GeneratesUUID(t *testing.T) {
	inst := &models.Instance{}
	err := inst.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, inst.ID)
	assert.Len(t, inst.ID, 36) // UUID format
}

func TestInstance_BeforeCreate_KeepsExistingID(t *testing.T) {
	inst := &models.Instance{ID: "custom-id"}
	err := inst.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.Equal(t, "custom-id", inst.ID)
}

func TestWebhook_BeforeCreate_GeneratesUUID(t *testing.T) {
	wh := &models.Webhook{}
	err := wh.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, wh.ID)
}

func TestMessage_BeforeCreate_GeneratesUUID(t *testing.T) {
	msg := &models.Message{}
	err := msg.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, msg.ID)
}

func TestLabel_BeforeCreate_GeneratesUUID(t *testing.T) {
	lbl := &models.Label{}
	err := lbl.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, lbl.ID)
}
