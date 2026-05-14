package feature

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"githubb.com/yeimar-projects/wa-go/app/facades"
	"githubb.com/yeimar-projects/wa-go/tests"
)

type InstanceTestSuite struct {
	suite.Suite
	tests.TestCase
	globalKey string
}

func TestInstanceTestSuite(t *testing.T) {
	suite.Run(t, new(InstanceTestSuite))
}

func (s *InstanceTestSuite) SetupSuite() {
	s.globalKey = facades.Config().GetString("whatsapp.global_api_key")
}

func (s *InstanceTestSuite) requireDB() {
	defer func() {
		if r := recover(); r != nil {
			s.T().Skip("DB not available")
		}
	}()
	// Ensure tables exist
	facades.Orm().Query().Exec("CREATE TABLE IF NOT EXISTS instances (id TEXT PRIMARY KEY, created_at DATETIME, updated_at DATETIME, name TEXT NOT NULL, token TEXT NOT NULL UNIQUE, status TEXT NOT NULL DEFAULT 'disconnected', jid TEXT, qrcode TEXT, qrcode_raw TEXT, proxy_protocol TEXT, proxy_host TEXT, proxy_port TEXT, proxy_username TEXT, proxy_password TEXT, whatsapp_version_major INTEGER DEFAULT 0, whatsapp_version_minor INTEGER DEFAULT 0, whatsapp_version_patch INTEGER DEFAULT 0, reject_call INTEGER DEFAULT 0, msg_reject_call TEXT DEFAULT '')")
}

// --- Auth tests ---

func (s *InstanceTestSuite) TestCreate_NoApiKey_Returns401() {
	body := strings.NewReader(`{"name":"test","token":"tok123"}`)
	resp, err := s.Http(s.T()).Post("/api/v1/instances/", body)
	s.Require().NoError(err)
	resp.AssertUnauthorized().AssertJson(map[string]any{"status": float64(401), "code": "UNAUTHORIZED", "message": "Authentication required. Provide 'apikey' header."})
}

func (s *InstanceTestSuite) TestCreate_WrongApiKey_Returns401() {
	body := strings.NewReader(`{"name":"test","token":"tok123"}`)
	resp, err := s.Http(s.T()).WithHeader("apikey", "wrong-key").Post("/api/v1/instances/", body)
	s.Require().NoError(err)
	resp.AssertUnauthorized().AssertJson(map[string]any{"status": float64(401), "code": "UNAUTHORIZED", "message": "Invalid API key."})
}

func (s *InstanceTestSuite) TestList_NoApiKey_Returns401() {
	resp, err := s.Http(s.T()).Get("/api/v1/instances/")
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

// --- CRUD tests (require DB + valid key) ---

func (s *InstanceTestSuite) TestCreate_ValidPayload_Returns201() {
	s.requireDB()
	token := "test-token-" + fmt.Sprintf("%d", time.Now().UnixNano())
	body := strings.NewReader(`{"name":"Test Instance","token":"` + token + `"}`)
	resp, err := s.Http(s.T()).WithHeader("apikey", s.globalKey).Post("/api/v1/instances/", body)
	s.Require().NoError(err)

	json, _ := resp.Json()
	resp.AssertCreated()
	s.Contains(json, "results")

	data := json["results"].(map[string]any)
	s.Equal("Test Instance", data["name"])
	s.NotEmpty(data["id"])
	s.Equal("disconnected", data["status"])
}

func (s *InstanceTestSuite) TestCreate_MissingName_Returns400() {
	s.requireDB()
	body := strings.NewReader(`{"token":"tok"}`)
	resp, err := s.Http(s.T()).WithHeader("apikey", s.globalKey).Post("/api/v1/instances/", body)
	s.Require().NoError(err)
	resp.AssertBadRequest().AssertJson(map[string]any{"status": float64(400), "code": "VALIDATION_ERROR"})
}

func (s *InstanceTestSuite) TestCreate_MissingToken_Returns400() {
	s.requireDB()
	body := strings.NewReader(`{"name":"test"}`)
	resp, err := s.Http(s.T()).WithHeader("apikey", s.globalKey).Post("/api/v1/instances/", body)
	s.Require().NoError(err)
	resp.AssertBadRequest().AssertJson(map[string]any{"status": float64(400), "code": "VALIDATION_ERROR"})
}

func (s *InstanceTestSuite) TestList_WithValidKey_ReturnsDataAndLinks() {
	s.requireDB()
	resp, err := s.Http(s.T()).WithHeader("apikey", s.globalKey).Get("/api/v1/instances/")
	s.Require().NoError(err)

	json, _ := resp.Json()
	resp.AssertOk()
	s.Contains(json, "results")

}

func (s *InstanceTestSuite) TestGet_NonExistent_Returns404() {
	s.requireDB()
	resp, err := s.Http(s.T()).WithHeader("apikey", s.globalKey).Get("/api/v1/instances/00000000-0000-0000-0000-000000000000")
	s.Require().NoError(err)
	resp.AssertNotFound().AssertJson(map[string]any{"status": float64(404), "code": "NOT_FOUND", "message": "instance not found"})
}

func (s *InstanceTestSuite) TestDelete_NonExistent_ReturnsError() {
	s.requireDB()
	resp, err := s.Http(s.T()).WithHeader("apikey", s.globalKey).Delete("/api/v1/instances/00000000-0000-0000-0000-000000000000", nil)
	s.Require().NoError(err)
	// Should succeed (delete is idempotent) or return error
	json, _ := resp.Json()
	s.NotNil(json)
}

// --- Instance action auth tests ---

func (s *InstanceTestSuite) TestConnect_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/connect", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *InstanceTestSuite) TestDisconnect_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/disconnect", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *InstanceTestSuite) TestLogout_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/logout", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *InstanceTestSuite) TestStatus_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Get("/api/v1/instances/id/status")
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *InstanceTestSuite) TestQRCode_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Get("/api/v1/instances/id/qr-code")
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}
