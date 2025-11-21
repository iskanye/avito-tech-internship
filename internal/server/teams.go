package server

import (
	"context"
	"errors"

	"github.com/iskanye/avito-tech-internship/internal/models"
	"github.com/iskanye/avito-tech-internship/internal/service/prassignment"
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
		teamReq.Members[i].UserID = member.UserId
		teamReq.Members[i].Username = member.Username
		teamReq.Members[i].IsActive = member.IsActive
	}

	team, err := s.assign.AddTeam(c, teamReq)
	if errors.Is(err, prassignment.ErrTeamExists) {
		response := api.PostTeamAdd400JSONResponse{}
		response.Error.Code = api.TEAMEXISTS
		response.Error.Message = err.Error()
		return response, nil
	} else if err != nil {
		return nil, err
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
		return response, nil
	} else if err != nil {
		return nil, err
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
		teamRes.Members[i].UserId = member.UserID
		teamRes.Members[i].Username = member.UserID
		teamRes.Members[i].IsActive = member.IsActive
	}

	return &teamRes
}
