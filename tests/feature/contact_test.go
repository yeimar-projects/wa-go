package feature

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yeimar-sandbox/wa-go/tests"
)

type ContactTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestContactTestSuite(t *testing.T) {
	suite.Run(t, new(ContactTestSuite))
}

// --- Auth Tests ---

func (s *ContactTestSuite) TestCheckContact_WithoutToken_ReturnsUnauthorized() {
	body := strings.NewReader(`{"phones":["5491155667788"]}`)
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/contacts/check", body)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *ContactTestSuite) TestGetContactInfo_WithoutToken_ReturnsUnauthorized() {
	resp, err := s.Http(s.T()).Get("/api/v1/instances/id/contacts/123@s.whatsapp.net")
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *ContactTestSuite) TestGetProfilePicture_WithoutToken_ReturnsUnauthorized() {
	resp, err := s.Http(s.T()).Get("/api/v1/instances/id/contacts/123@s.whatsapp.net/profile-picture")
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *ContactTestSuite) TestGetBusinessProfile_WithoutToken_ReturnsUnauthorized() {
	resp, err := s.Http(s.T()).Get("/api/v1/instances/id/contacts/123@s.whatsapp.net/business-profile")
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *ContactTestSuite) TestBlockContact_WithoutToken_ReturnsUnauthorized() {
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/contacts/123@s.whatsapp.net/block", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *ContactTestSuite) TestUnblockContact_WithoutToken_ReturnsUnauthorized() {
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/contacts/123@s.whatsapp.net/unblock", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *ContactTestSuite) TestGetBlocklist_WithoutToken_ReturnsUnauthorized() {
	resp, err := s.Http(s.T()).Get("/api/v1/instances/id/contacts/blocklist")
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

// --- Payload Validation & Logic Tests ---

func (s *ContactTestSuite) TestCheckContact_MissingPhones_ReturnsValidationError() {
	inst := s.CreateTestInstance("disconnected")
	defer s.ClearDB()

	body := strings.NewReader(`{}`)
	resp, err := s.Http(s.T()).WithHeader("apikey", inst.Token).Post("/api/v1/instances/"+inst.ID+"/contacts/check", body)
	
	s.Require().NoError(err)
	resp.AssertBadRequest().AssertJson(map[string]any{"status": float64(400), "code": "VALIDATION_ERROR", "message": "'phones' array is required"})
}


