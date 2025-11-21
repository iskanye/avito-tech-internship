package server

import (
	"context"
	"errors"

	"github.com/iskanye/avito-tech-internship/internal/models"
	prassignment "github.com/iskanye/avito-tech-internship/internal/service/pr_assignment"
	"github.com/iskanye/avito-tech-internship/pkg/api"
)

// (POST /team/add)
func (s *serverAPI) PostTeamAdd(
	c context.Context,
	req api.PostTeamAddRequestObject,
) (api.PostTeamAddResponseObject, error) {
	teamReq := models.Team{
		TeamName: req.Body.TeamName,
		Members:  make([]models.User, len(req.Body.Members)),
	}
	for i, member := range req.Body.Members {
		teamMember := models.User{
			UserID:   member.UserId,
			Username: member.Username,
			IsActive: member.IsActive,
		}
		teamReq.Members[i] = teamMember
	}

	team, err := s.assign.AddTeam(c, teamReq)
	if errors.Is(err, prassignment.ErrTeamExists) {
		response := api.PostTeamAdd400JSONResponse{}
		response.Error.Code = api.TEAMEXISTS
		response.Error.Message = err.Error()
		return response, err
	}

	teamResp := convertTeamToApi(&team)
	response := api.PostTeamAdd201JSONResponse{
		Team: teamResp,
	}
	return response, nil
}

// (GET /team/get)
func (s *serverAPI) GetTeamGet(
	c context.Context,
	req api.GetTeamGetRequestObject,
) (api.GetTeamGetResponseObject, error) {
	team, err := s.assign.GetTeam(c, req.Params.TeamName)
	if errors.Is(err, prassignment.ErrNotFound) {
		response := api.GetTeamGet404JSONResponse{}
		response.Error.Code = api.NOTFOUND
		response.Error.Message = err.Error()
		return response, err
	}

	teamResp := convertTeamToApi(&team)
	response := (api.GetTeamGet200JSONResponse)(*teamResp)
	return response, nil
}

func convertTeamToApi(team *models.Team) *api.Team {
	teamRes := api.Team{
		TeamName: team.TeamName,
		Members:  make([]api.TeamMember, len(team.Members)),
	}
	for i, member := range team.Members {
		teamMember := api.TeamMember{
			UserId:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		}
		teamRes.Members[i] = teamMember
	}

	return &teamRes
}
