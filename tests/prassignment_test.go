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

	NOT_FOUND   = "resource not found"
	TEAM_EXISTS = "team_name already exists"
)

func TestTeams_AddGetTeam_Success(t *testing.T) {
	s, ctx := suite.New(t)

	team := suite.RandomTeam(membersCount)

	// Добавить команду
	addTeamResp, err := s.Client.PostTeamAddWithResponse(ctx, *team)
	require.NoError(t, err)
	require.NotEmpty(t, addTeamResp.JSON201)
	suite.RequireTeamsEqual(t, team, addTeamResp.JSON201.Team)

	// Получить команду
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

	// Добавляем две одинаковые команды
	resp, err := s.Client.PostTeamAddWithResponse(ctx, *team)
	require.NoError(t, err)
	require.NotEmpty(t, resp.JSON201)
	suite.RequireTeamsEqual(t, team, resp.JSON201.Team)

	resp, err = s.Client.PostTeamAddWithResponse(ctx, *team)
	require.NoError(t, err)
	require.NotEmpty(t, resp.JSON400)
	assert.Equal(t, api.TEAMEXISTS, resp.JSON400.Error.Code)
	assert.Equal(t, TEAM_EXISTS, resp.JSON400.Error.Message)
}

func TestTeam_GetTeam_NotFound(t *testing.T) {
	s, ctx := suite.New(t)

	// Получить команду, которая не существует
	resp, err := s.Client.GetTeamGetWithResponse(ctx, &api.GetTeamGetParams{
		TeamName: gofakeit.Noun(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.JSON404)
	assert.Equal(t, api.NOTFOUND, resp.JSON404.Error.Code)
	assert.Equal(t, NOT_FOUND, resp.JSON404.Error.Message)
}

func TestUser_SetIsActive_Success(t *testing.T) {
	s, ctx := suite.New(t)

	team := suite.RandomTeam(membersCount)

	// Создаем команду
	addTeam, err := s.Client.PostTeamAddWithResponse(ctx, *team)
	require.NoError(t, err)
	require.NotEmpty(t, addTeam.JSON201)
	suite.RequireTeamsEqual(t, team, addTeam.JSON201.Team)

	// Изменяем состояние активности одного пользователя
	setIsActiveReq := api.PostUsersSetIsActiveJSONRequestBody{
		UserId:   team.Members[0].UserId,
		IsActive: !team.Members[0].IsActive,
	}

	setIsActive, err := s.Client.PostUsersSetIsActiveWithResponse(ctx, setIsActiveReq)
	require.NoError(t, err)
	require.NotEmpty(t, setIsActive.JSON200)

	assert.Equal(t, team.Members[0].UserId, setIsActive.JSON200.User.UserId)
	assert.Equal(t, team.Members[0].Username, setIsActive.JSON200.User.Username)
	assert.Equal(t, team.TeamName, setIsActive.JSON200.User.TeamName)
	assert.Equal(t, team.Members[0].IsActive, !setIsActive.JSON200.User.IsActive)
}

func TestUser_SetIsActive_NotFound(t *testing.T) {
	s, ctx := suite.New(t)

	// Пытаемся изменить состояние активности несуществующего пользователя
	req := api.PostUsersSetIsActiveJSONRequestBody{
		UserId:   gofakeit.UUID(),
		IsActive: false,
	}

	resp, err := s.Client.PostUsersSetIsActiveWithResponse(ctx, req)
	require.NoError(t, err)
	require.NotEmpty(t, resp.JSON404)
	assert.Equal(t, api.NOTFOUND, resp.JSON404.Error.Code)
	assert.Equal(t, NOT_FOUND, resp.JSON404.Error.Message)
}
