package feature

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"githubb.com/yeimar-projects/wa-go/tests"
)

type GroupTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestGroupTestSuite(t *testing.T) {
	suite.Run(t, new(GroupTestSuite))
}

// --- Auth Tests ---

func (s *GroupTestSuite) TestListGroups_WithoutToken_ReturnsUnauthorized() {
	resp, err := s.Http(s.T()).Get("/api/v1/instances/id/groups")
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *GroupTestSuite) TestCreateGroup_WithoutToken_ReturnsUnauthorized() {
	body := strings.NewReader(`{"subject":"My Group","participants":["5491155667788"]}`)
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/groups", body)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *GroupTestSuite) TestGetGroupInfo_WithoutToken_ReturnsUnauthorized() {
	resp, err := s.Http(s.T()).Get("/api/v1/instances/id/groups/120363001234567890@g.us")
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *GroupTestSuite) TestJoinGroupWithLink_WithoutToken_ReturnsUnauthorized() {
	body := strings.NewReader(`{"inviteCode":"AbCdEfGhIjK"}`)
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/groups/g1/join", body)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *GroupTestSuite) TestLeaveGroup_WithoutToken_ReturnsUnauthorized() {
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/groups/120363001234567890@g.us/leave", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *GroupTestSuite) TestGetGroupInviteLink_WithoutToken_ReturnsUnauthorized() {
	resp, err := s.Http(s.T()).Get("/api/v1/instances/id/groups/120363001234567890@g.us/invite-link")
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *GroupTestSuite) TestResetGroupInviteLink_WithoutToken_ReturnsUnauthorized() {
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/groups/g1/invite-link/reset", nil)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *GroupTestSuite) TestAddParticipants_WithoutToken_ReturnsUnauthorized() {
	body := strings.NewReader(`{"participants":["5491155667788","5491122334455"]}`)
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/groups/g1/participants/add", body)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *GroupTestSuite) TestRemoveParticipants_WithoutToken_ReturnsUnauthorized() {
	body := strings.NewReader(`{"participants":["5491155667788"]}`)
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/groups/g1/participants/remove", body)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *GroupTestSuite) TestPromoteParticipants_WithoutToken_ReturnsUnauthorized() {
	body := strings.NewReader(`{"participants":["5491155667788"]}`)
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/groups/g1/participants/promote", body)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *GroupTestSuite) TestDemoteParticipants_WithoutToken_ReturnsUnauthorized() {
	body := strings.NewReader(`{"participants":["5491155667788"]}`)
	resp, err := s.Http(s.T()).Post("/api/v1/instances/id/groups/g1/participants/demote", body)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

func (s *GroupTestSuite) TestUpdateGroupSettings_WithoutToken_ReturnsUnauthorized() {
	body := strings.NewReader(`{"name":"New Group Name","description":"Updated desc"}`)
	resp, err := s.Http(s.T()).Patch("/api/v1/instances/id/groups/g1/settings", body)
	s.Require().NoError(err)
	resp.AssertUnauthorized()
}

// --- Payload Validation & Logic Tests ---

func (s *GroupTestSuite) TestCreateGroup_MissingSubject_ReturnsValidationError() {
	inst := s.CreateTestInstance("disconnected")
	defer s.ClearDB()

	body := strings.NewReader(`{"participants":["5491155667788"]}`)
	resp, err := s.Http(s.T()).WithHeader("apikey", inst.Token).Post("/api/v1/instances/"+inst.ID+"/groups", body)
	
	s.Require().NoError(err)
	resp.AssertBadRequest().AssertJson(map[string]any{"status": float64(400), "code": "VALIDATION_ERROR", "message": "'subject' is required"})
}

func (s *GroupTestSuite) TestAddParticipants_MissingParticipants_ReturnsValidationError() {
	inst := s.CreateTestInstance("disconnected")
	defer s.ClearDB()

	body := strings.NewReader(`{}`)
	resp, err := s.Http(s.T()).WithHeader("apikey", inst.Token).Post("/api/v1/instances/"+inst.ID+"/groups/123@g.us/participants/add", body)
	
	s.Require().NoError(err)
	resp.AssertBadRequest().AssertJson(map[string]any{"status": float64(400), "code": "VALIDATION_ERROR", "message": "'participants' array is required"})
}

func (s *GroupTestSuite) TestListGroups_WithValidTokenButDisconnected_ReturnsError() {
	inst := s.CreateTestInstance("disconnected")
	defer s.ClearDB()

	resp, err := s.Http(s.T()).WithHeader("apikey", inst.Token).Get("/api/v1/instances/"+inst.ID+"/groups")
	
	s.Require().NoError(err)
	resp.AssertJson(map[string]any{"status": float64(503), "code": "WHATSAPP_NOT_CONNECTED"})
}
