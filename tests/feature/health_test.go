package feature

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"githubb.com/yeimar-projects/wa-go/tests"
)

type HealthTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestHealthTestSuite(t *testing.T) {
	suite.Run(t, new(HealthTestSuite))
}

func (s *HealthTestSuite) TestHealth_ReturnsOk() {
	resp, err := s.Http(s.T()).Get("/api/v1/health")
	s.Require().NoError(err)
	resp.AssertOk().AssertJson(map[string]any{"status": "ok"})
}

func (s *HealthTestSuite) TestHealth_ContentTypeIsJSON() {
	resp, err := s.Http(s.T()).Get("/api/v1/health")
	s.Require().NoError(err)
	resp.AssertOk().AssertHeader("Content-Type", "application/json; charset=utf-8")
}
