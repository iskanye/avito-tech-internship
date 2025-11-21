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

func TestTeams_AddTeam_Success(t *testing.T) {
	s, ctx := suite.New(t)

	team := s.RandomTeam(membersCount)

	resp, err := s.Client.PostTeamAddWithResponse(ctx, *team)
	require.NoError(t, err)
	require.IsType(t, &api.PostTeamAdd201JSONResponse{}, resp)
	require.Equal(t, *team, *resp.JSON201.Team)
}
