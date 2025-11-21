package tests

import (
	"testing"

	"github.com/iskanye/avito-tech-internship/pkg/api"
	"github.com/iskanye/avito-tech-internship/tests/suite"
	"github.com/stretchr/testify/require"
)

const (
	membersCount = 5
)

func TestTeams_AddGetTeam_Success(t *testing.T) {
	s, ctx := suite.New(t)

	team := s.RandomTeam(membersCount)

	addTeamResp, err := s.Client.PostTeamAddWithResponse(ctx, *team)
	require.NoError(t, err)
	require.NotEmpty(t, addTeamResp.JSON201)
	suite.RequireTeamsEqual(t, team, addTeamResp.JSON201.Team)

	getTeamResp, err := s.Client.GetTeamGetWithResponse(ctx, &api.GetTeamGetParams{
		TeamName: team.TeamName,
	})
	require.NoError(t, err)
	require.NotEmpty(t, getTeamResp.JSON200)
	suite.RequireTeamsEqual(t, team, getTeamResp.JSON200)
}
