package feature

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yeimar-sandbox/wa-go/tests"
)

type WebhookTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestWebhookTestSuite(t *testing.T) {
	suite.Run(t, new(WebhookTestSuite))
}

// --- Auth Tests ---

func (s *WebhookTestSuite) TestCreateWebhook_WithoutToken_ReturnsUnauthorized() {
	body := strings.NewReader(`{"url":"https://example.com/hook","events":["message.received"]}`)
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/webhooks", body)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *WebhookTestSuite) TestListWebhooks_WithoutToken_ReturnsUnauthorized() {
	resp, err := s.Http(s.T()).Get("/api/v1/instances/id/webhooks")
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *WebhookTestSuite) TestGetWebhookInfo_WithoutToken_ReturnsUnauthorized() {
	resp, err := s.Http(s.T()).Get("/api/v1/instances/id/webhooks/wh1")
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *WebhookTestSuite) TestDeleteWebhook_WithoutToken_ReturnsUnauthorized() {
	resp, err := s.Http(s.T()).Delete("/api/v1/instances/id/webhooks/wh1", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *WebhookTestSuite) TestTriggerWebhookTest_WithoutToken_ReturnsUnauthorized() {
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/webhooks/wh1/test", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

// --- Payload Validation Tests ---

func (s *WebhookTestSuite) TestCreateWebhook_MissingUrl_ReturnsValidationError() {
	inst := s.CreateTestInstance("disconnected")
	defer s.ClearDB()

	body := strings.NewReader(`{"events":["message.received"]}`)
	resp, err := s.Http(s.T()).WithHeader("apikey", inst.Token).Post("/api/v1/instances/"+inst.ID+"/webhooks", body)
	
	s.Require().NoError(err)
	resp.AssertBadRequest().AssertJson(map[string]any{"status": float64(400), "code": "VALIDATION_ERROR", "message": "'url' is required"})
}

func (s *WebhookTestSuite) TestCreateWebhook_MissingEvents_ReturnsValidationError() {
	inst := s.CreateTestInstance("disconnected")
	defer s.ClearDB()

	body := strings.NewReader(`{"url":"https://example.com/hook"}`)
	resp, err := s.Http(s.T()).WithHeader("apikey", inst.Token).Post("/api/v1/instances/"+inst.ID+"/webhooks", body)
	
	s.Require().NoError(err)
	resp.AssertBadRequest().AssertJson(map[string]any{"status": float64(400), "code": "VALIDATION_ERROR", "message": "'events' array is required and must not be empty"})
}
