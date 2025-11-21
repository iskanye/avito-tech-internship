package tests

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/iskanye/avito-tech-internship/pkg/api"
	"github.com/iskanye/avito-tech-internship/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	membersCount = 5

	resourseNotFoundMsg = "resource not found"
	teamExistsMsg       = "team_name already exists"
)

func TestTeams_AddGetTeam_Success(t *testing.T) {
	s, ctx := suite.New(t)

	team := suite.RandomTeam(membersCount)

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

func TestTeam_AddTeam_Dublicate(t *testing.T) {
	s, ctx := suite.New(t)

	team := suite.RandomTeam(membersCount)

	resp, err := s.Client.PostTeamAddWithResponse(ctx, *team)
	require.NoError(t, err)
	require.NotEmpty(t, resp.JSON201)
	suite.RequireTeamsEqual(t, team, resp.JSON201.Team)

	resp, err = s.Client.PostTeamAddWithResponse(ctx, *team)
	require.NoError(t, err)
	require.NotEmpty(t, resp.JSON400)
	assert.Equal(t, api.TEAMEXISTS, resp.JSON400.Error.Code)
	assert.Equal(t, teamExistsMsg, resp.JSON400.Error.Message)
}

func TestTeam_GetTeam_NotFound(t *testing.T) {
	s, ctx := suite.New(t)

	resp, err := s.Client.GetTeamGetWithResponse(ctx, &api.GetTeamGetParams{
		TeamName: gofakeit.Noun(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.JSON404)
	assert.Equal(t, api.NOTFOUND, resp.JSON404.Error.Code)
	assert.Equal(t, resourseNotFoundMsg, resp.JSON404.Error.Message)
}
