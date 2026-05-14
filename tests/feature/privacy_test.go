package feature

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"githubb.com/yeimar-projects/wa-go/tests"
)

type PrivacyTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestPrivacyTestSuite(t *testing.T) {
	suite.Run(t, new(PrivacyTestSuite))
}

func (s *PrivacyTestSuite) TestGet_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Get("/api/v1/instances/id/privacy")
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *PrivacyTestSuite) TestUpdate_NoToken_Returns401() {
	body := strings.NewReader(`{"setting":"last_seen","value":"nobody"}`)
	resp, err := s.Http(s.T()).Patch("/api/v1/instances/id/privacy", body)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}
