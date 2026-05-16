package feature

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yeimar-sandbox/wa-go/tests"
)

type CallTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestCallTestSuite(t *testing.T) {
	suite.Run(t, new(CallTestSuite))
}

func (s *CallTestSuite) TestReject_NoToken_Returns401() {
	body := strings.NewReader(`{"callCreator":"123@s.whatsapp.net","callId":"call1"}`)
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/calls/call1/reject", body)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}
