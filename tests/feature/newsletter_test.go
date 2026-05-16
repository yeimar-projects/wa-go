package feature

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yeimar-sandbox/wa-go/tests"
)

type NewsletterTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestNewsletterTestSuite(t *testing.T) {
	suite.Run(t, new(NewsletterTestSuite))
}

func (s *NewsletterTestSuite) TestList_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Get("/api/v1/instances/id/newsletters")
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *NewsletterTestSuite) TestGet_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Get("/api/v1/instances/id/newsletters/nl1")
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *NewsletterTestSuite) TestFollow_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/newsletters/nl1/follow", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *NewsletterTestSuite) TestUnfollow_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/newsletters/nl1/unfollow", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *NewsletterTestSuite) TestMute_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/newsletters/nl1/mute", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *NewsletterTestSuite) TestUnmute_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/newsletters/nl1/unmute", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}
