package feature

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"githubb.com/yeimar-projects/wa-go/tests"
)

type MessageTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestMessageTestSuite(t *testing.T) {
	suite.Run(t, new(MessageTestSuite))
}

// --- Auth Tests ---

func (s *MessageTestSuite) TestSendText_WithoutToken_ReturnsUnauthorized() {
	body := strings.NewReader(`{"to":"123","type":"text","text":{"body":"hi"}}`)
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/messages", body)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *MessageTestSuite) TestReactToMessage_WithoutToken_ReturnsUnauthorized() {
	body := strings.NewReader(`{"emoji":"👍"}`)
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/messages/msg1/react", body)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *MessageTestSuite) TestRevokeMessage_WithoutToken_ReturnsUnauthorized() {
	body := strings.NewReader(`{"chatJid":"123@s.whatsapp.net"}`)
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/messages/msg1/revoke", body)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *MessageTestSuite) TestEditMessage_WithoutToken_ReturnsUnauthorized() {
	body := strings.NewReader(`{"chatJid":"123@s.whatsapp.net","text":"new"}`)
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/messages/msg1/edit", body)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *MessageTestSuite) TestMarkMessageAsRead_WithoutToken_ReturnsUnauthorized() {
	body := strings.NewReader(`{"chatJid":"123@s.whatsapp.net","messageIds":["m1"]}`)
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/messages/msg1/read", body)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

// --- Payload Validation & Logic Tests ---



func (s *MessageTestSuite) TestSendMessage_MissingType_ReturnsValidationError() {
	inst := s.CreateTestInstance("disconnected")
	defer s.ClearDB()

	body := strings.NewReader(`{"to":"5491155667788","text":{"body":"Missing type"}}`)
	resp, err := s.Http(s.T()).WithHeader("apikey", inst.Token).Post("/api/v1/instances/"+inst.ID+"/messages", body)
	
	s.Require().NoError(err)
	resp.AssertBadRequest().AssertJson(map[string]any{"status": float64(400), "code": "VALIDATION_ERROR", "message": "'type' field is required"})
}

func (s *MessageTestSuite) TestSendMessage_MissingTo_ReturnsValidationError() {
	inst := s.CreateTestInstance("disconnected")
	defer s.ClearDB()

	body := strings.NewReader(`{"type":"text","text":{"body":"Missing type"}}`)
	resp, err := s.Http(s.T()).WithHeader("apikey", inst.Token).Post("/api/v1/instances/"+inst.ID+"/messages", body)
	
	s.Require().NoError(err)
	resp.AssertBadRequest().AssertJson(map[string]any{"status": float64(400), "code": "VALIDATION_ERROR", "message": "'to' field is required"})
}



func (s *MessageTestSuite) TestReactToMessage_MissingEmoji_ReturnsValidationError() {
	inst := s.CreateTestInstance("disconnected")
	defer s.ClearDB()

	body := strings.NewReader(`{"wrongField":"👍"}`)
	resp, err := s.Http(s.T()).WithHeader("apikey", inst.Token).Post("/api/v1/instances/"+inst.ID+"/messages/msg1/react", body)
	
	s.Require().NoError(err)
	resp.AssertBadRequest().AssertJson(map[string]any{"status": float64(400), "code": "VALIDATION_ERROR", "message": "'emoji' is required"})
}
