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
	PR_EXISTS   = "PR id already exists"
)

func TestTeams_AddGetTeam_Success(t *testing.T) {
	s, ctx := suite.New(t)

	team := suite.RandomTeam(membersCount, gofakeit.Bool)

	// Добавить команду
	addTeamResp, err := s.Client.PostTeamAddWithResponse(ctx, *team)
	require.NoError(t, err)
	require.NotEmpty(t, addTeamResp.JSON201)
	suite.CheckTeamsEqual(t, team, addTeamResp.JSON201.Team)

	// Получить команду
	getTeamResp, err := s.Client.GetTeamGetWithResponse(ctx, &api.GetTeamGetParams{
		TeamName: team.TeamName,
	})
	require.NoError(t, err)
	require.NotEmpty(t, getTeamResp.JSON200)
	suite.CheckTeamsEqual(t, team, getTeamResp.JSON200)
}

func TestTeams_AddTeam_Dublicate(t *testing.T) {
	s, ctx := suite.New(t)

	team := suite.RandomTeam(membersCount, gofakeit.Bool)

	// Добавляем две одинаковые команды
	resp, err := s.Client.PostTeamAddWithResponse(ctx, *team)
	require.NoError(t, err)
	require.NotEmpty(t, resp.JSON201)
	suite.CheckTeamsEqual(t, team, resp.JSON201.Team)

	resp, err = s.Client.PostTeamAddWithResponse(ctx, *team)
	require.NoError(t, err)
	require.NotEmpty(t, resp.JSON400)
	assert.Equal(t, api.TEAMEXISTS, resp.JSON400.Error.Code)
	assert.Equal(t, TEAM_EXISTS, resp.JSON400.Error.Message)
}

func TestTeams_GetTeam_NotFound(t *testing.T) {
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

func TestUsers_SetIsActive_Success(t *testing.T) {
	s, ctx := suite.New(t)

	team := suite.RandomTeam(membersCount, gofakeit.Bool)

	// Создаем команду
	addTeam, err := s.Client.PostTeamAddWithResponse(ctx, *team)
	require.NoError(t, err)
	require.NotEmpty(t, addTeam.JSON201)
	suite.CheckTeamsEqual(t, team, addTeam.JSON201.Team)

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

func TestUsers_SetIsActive_NotFound(t *testing.T) {
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

func TestUsers_GetReview_Success(t *testing.T) {
	s, ctx := suite.New(t)

	// Создаем команду из 3 активных человек, чтобы все учавствовали в пул реквесте
	team := suite.RandomTeam(3, func() bool { return true })

	// Добавить команду
	addTeamResp, err := s.Client.PostTeamAddWithResponse(ctx, *team)
	require.NoError(t, err)
	require.NotEmpty(t, addTeamResp.JSON201)
	suite.CheckTeamsEqual(t, team, addTeamResp.JSON201.Team)

	// Cоздаем два пул реквеста с одним автором
	pr1 := suite.RandomPullRequest(team.Members[0].UserId)

	addPullRequest, err := s.Client.PostPullRequestCreateWithResponse(ctx, *pr1)
	require.NoError(t, err)
	require.NotEmpty(t, addPullRequest.JSON201)
	suite.CheckPullRequestEqual(t, pr1, addPullRequest.JSON201.Pr)

	pr2 := suite.RandomPullRequest(team.Members[0].UserId)

	addPullRequest, err = s.Client.PostPullRequestCreateWithResponse(ctx, *pr2)
	require.NoError(t, err)
	require.NotEmpty(t, addPullRequest.JSON201)
	suite.CheckPullRequestEqual(t, pr2, addPullRequest.JSON201.Pr)

	// Получаем список пул реквестов второго юзера
	getReview, err := s.Client.GetUsersGetReviewWithResponse(ctx, &api.GetUsersGetReviewParams{
		UserId: team.Members[1].UserId,
	})
	require.NoError(t, err)
	require.NotEmpty(t, getReview.JSON200)
	suite.CheckPullRequestsEqual(t,
		getReview.JSON200.PullRequests,
		suite.PullRequestCreateToModel(pr1),
		suite.PullRequestCreateToModel(pr2),
	)
}

func TestPullRequests_CreatePullRequest_Success(t *testing.T) {
	s, ctx := suite.New(t)

	team := suite.RandomTeam(membersCount, gofakeit.Bool)

	// Добавить команду
	addTeam, err := s.Client.PostTeamAddWithResponse(ctx, *team)
	require.NoError(t, err)
	require.NotEmpty(t, addTeam.JSON201)
	suite.CheckTeamsEqual(t, team, addTeam.JSON201.Team)

	pullRequest := suite.RandomPullRequest(team.Members[0].UserId)

	// Добавляем пул реквест
	addPullRequest, err := s.Client.PostPullRequestCreateWithResponse(ctx, *pullRequest)
	require.NoError(t, err)
	require.NotEmpty(t, addPullRequest.JSON201)
	suite.CheckPullRequestEqual(t, pullRequest, addPullRequest.JSON201.Pr)
}

func TestPullRequests_CreatePullRequest_Dublicate(t *testing.T) {
	s, ctx := suite.New(t)

	team := suite.RandomTeam(membersCount, gofakeit.Bool)

	// Добавить команду
	addTeam, err := s.Client.PostTeamAddWithResponse(ctx, *team)
	require.NoError(t, err)
	require.NotEmpty(t, addTeam.JSON201)
	suite.CheckTeamsEqual(t, team, addTeam.JSON201.Team)

	pullRequest := suite.RandomPullRequest(team.Members[0].UserId)

	// Добавляем пул реквест
	addPullRequest, err := s.Client.PostPullRequestCreateWithResponse(ctx, *pullRequest)
	require.NoError(t, err)
	require.NotEmpty(t, addPullRequest.JSON201)
	suite.CheckPullRequestEqual(t, pullRequest, addPullRequest.JSON201.Pr)

	// Вставляем дубликат
	addPullRequest, err = s.Client.PostPullRequestCreateWithResponse(ctx, *pullRequest)
	require.NoError(t, err)
	require.NotEmpty(t, addPullRequest.JSON409)
	assert.Equal(t, api.PREXISTS, addPullRequest.JSON409.Error.Code)
	assert.Equal(t, PR_EXISTS, addPullRequest.JSON409.Error.Message)
}

func TestPullRequests_CreatePullRequest_NotFound(t *testing.T) {
	s, ctx := suite.New(t)

	pullRequest := suite.RandomPullRequest(gofakeit.UUID())

	// Добавляем пул реквест с несуществующим автором
	addPullRequest, err := s.Client.PostPullRequestCreateWithResponse(ctx, *pullRequest)
	require.NoError(t, err)
	require.NotEmpty(t, addPullRequest.JSON404)
	assert.Equal(t, api.NOTFOUND, addPullRequest.JSON404.Error.Code)
	assert.Equal(t, NOT_FOUND, addPullRequest.JSON404.Error.Message)
}
