package feature

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yeimar-sandbox/wa-go/tests"
)

type PresenceTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestPresenceTestSuite(t *testing.T) {
	suite.Run(t, new(PresenceTestSuite))
}

func (s *PresenceTestSuite) TestSet_NoToken_Returns401() {
	body := strings.NewReader(`{"presence":"available"}`)
	resp, err := s.Http(s.T()).Put("/api/v1/instances/id/presence", body)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *PresenceTestSuite) TestSubscribe_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/presence/123@s.whatsapp.net/subscribe", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}
