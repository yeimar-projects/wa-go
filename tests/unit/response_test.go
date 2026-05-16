package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yeimar-sandbox/wa-go/app/http/response"
)

func TestNewSuccess(t *testing.T) {
	data := map[string]string{"id": "123"}
	r := response.NewSuccess(data, "ok")

	assert.Equal(t, 200, r.Status)
	assert.Equal(t, "SUCCESS", r.Code)
	assert.Equal(t, "ok", r.Message)
	assert.Equal(t, data, r.Results)
}

func TestNewPaginated(t *testing.T) {
	r := response.NewPaginated([]string{"a", "b"}, 20, true, "cursor_abc", "", "ok")

	assert.Equal(t, 20, r.Pagination.Limit)
	assert.True(t, r.Pagination.HasMore)
	assert.Equal(t, "cursor_abc", r.Pagination.Cursors.After)
}

func TestNewError(t *testing.T) {
	r := response.NewError(404, "NOT_FOUND", "resource not found")

	assert.Equal(t, 404, r.Status)
	assert.Equal(t, "NOT_FOUND", r.Code)
	assert.Equal(t, "resource not found", r.Message)
}

func TestInstanceLinks_Connected(t *testing.T) {
	links := response.InstanceLinks("/api/v1", "inst_123", "connected")

	assert.Contains(t, links, "self")
	assert.Contains(t, links, "disconnect")
	assert.Contains(t, links, "messages")
	assert.Contains(t, links, "groups")
	assert.Contains(t, links, "contacts")
	assert.Contains(t, links, "newsletters")
	assert.NotContains(t, links, "connect")
	assert.Equal(t, "/api/v1/instances/inst_123/disconnect", links["disconnect"].Href)
	assert.Equal(t, "POST", links["disconnect"].Method)
}

func TestInstanceLinks_Disconnected(t *testing.T) {
	links := response.InstanceLinks("/api/v1", "inst_123", "disconnected")

	assert.Contains(t, links, "connect")
	assert.Contains(t, links, "pair-phone")
	assert.Contains(t, links, "delete")
	assert.NotContains(t, links, "disconnect")
	assert.NotContains(t, links, "messages")
}

func TestInstanceLinks_QRCode(t *testing.T) {
	links := response.InstanceLinks("/api/v1", "inst_123", "qr_code")

	assert.Contains(t, links, "qr-code")
	assert.Contains(t, links, "disconnect")
	assert.NotContains(t, links, "connect")
}
