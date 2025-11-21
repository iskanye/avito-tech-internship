package suite

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/iskanye/avito-tech-internship/pkg/api"
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
