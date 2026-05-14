package feature

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"githubb.com/yeimar-projects/wa-go/tests"
)

type ProfileTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestProfileTestSuite(t *testing.T) {
	suite.Run(t, new(ProfileTestSuite))
}

func (s *ProfileTestSuite) TestSetStatus_NoToken_Returns401() {
	body := strings.NewReader(`{"message":"Hello world"}`)
	resp, err := s.Http(s.T()).Put("/api/v1/instances/id/profile/status-message", body)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *ProfileTestSuite) TestGetQRLink_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Get("/api/v1/instances/id/profile/qr-link")
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *ProfileTestSuite) TestRevokeQRLink_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/profile/qr-link/revoke", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}
