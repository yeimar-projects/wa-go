package feature

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yeimar-sandbox/wa-go/tests"
)

type ChatTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestChatTestSuite(t *testing.T) {
	suite.Run(t, new(ChatTestSuite))
}

func (s *ChatTestSuite) TestPin_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/chats/chat1/pin", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *ChatTestSuite) TestUnpin_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/chats/chat1/unpin", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *ChatTestSuite) TestArchive_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/chats/chat1/archive", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *ChatTestSuite) TestUnarchive_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/chats/chat1/unarchive", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *ChatTestSuite) TestMute_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/chats/chat1/mute", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *ChatTestSuite) TestUnmute_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/chats/chat1/unmute", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *ChatTestSuite) TestPresence_NoToken_Returns401() {
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/chats/chat1/presence", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}
