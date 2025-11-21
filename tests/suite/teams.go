package suite

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/iskanye/avito-tech-internship/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *Suite) RandomTeam(membersCount int) *api.Team {
	team := &api.Team{
		TeamName: gofakeit.Noun(),
		Members:  make([]api.TeamMember, membersCount),
	}

	for i := range membersCount {
		member := api.TeamMember{
			IsActive: gofakeit.Bool(),
			UserId:   gofakeit.UUID(),
			Username: gofakeit.Username(),
		}
		team.Members[i] = member
	}

	return team
}

func RequireTeamsEqual(t *testing.T, team1 *api.Team, team2 *api.Team) {
	require.Equal(t, team1.TeamName, team2.TeamName)

	membersSet := make(map[string]teamMember)
	for _, member := range team1.Members {
		membersSet[member.UserId] = teamMember{
			username: member.Username,
			isActive: member.IsActive,
		}
	}

	for _, member1 := range team2.Members {
		member2, ok := membersSet[member1.UserId]
		require.True(t, ok)
		assert.Equal(t, member1.Username, member2.username)
		assert.Equal(t, member1.IsActive, member2.isActive)
	}
}

type teamMember struct {
	username string
	isActive bool
}
