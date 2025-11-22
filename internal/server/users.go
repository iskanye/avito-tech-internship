package server

import (
	"context"
	"errors"

	"github.com/iskanye/avito-tech-internship/internal/service/prassignment"
	"github.com/iskanye/avito-tech-internship/pkg/api"
)

// (GET /users/getReview)
func (s *serverAPI) GetUsersGetReview(
	c context.Context,
	req api.GetUsersGetReviewRequestObject,
) (api.GetUsersGetReviewResponseObject, error) {
	pullRequests, err := s.assign.GetReview(c, req.Params.UserId)
	if errors.Is(err, prassignment.ErrNotFound) {
		response := api.GetUsersGetReview404JSONResponse{}
		response.Error.Code = api.NOTFOUND
		response.Error.Message = err.Error()
		return response, nil
	}
	if err != nil {
		return nil, err
	}

	response := api.GetUsersGetReview200JSONResponse{
		PullRequests: make([]api.PullRequestShort, len(pullRequests)),
		UserId:       req.Params.UserId,
	}
	for i, pullRequest := range pullRequests {
		response.PullRequests[i].PullRequestId = pullRequest.ID
		response.PullRequests[i].PullRequestName = pullRequest.Name
		response.PullRequests[i].AuthorId = pullRequest.AuthorID
		response.PullRequests[i].Status = api.PullRequestShortStatus(pullRequest.Status)
	}

	return response, nil
}

// (POST /users/setIsActive)
func (s *serverAPI) PostUsersSetIsActive(
	c context.Context,
	req api.PostUsersSetIsActiveRequestObject,
) (api.PostUsersSetIsActiveResponseObject, error) {
	user, err := s.assign.SetIsActive(c, req.Body.UserId, req.Body.IsActive)
	if errors.Is(err, prassignment.ErrNotFound) {
		response := api.PostUsersSetIsActive404JSONResponse{}
		response.Error.Code = api.NOTFOUND
		response.Error.Message = err.Error()
		return response, nil
	}
	if err != nil {
		return nil, err
	}

	response := api.PostUsersSetIsActive200JSONResponse{}
	response.User = &api.User{
		UserId:   user.UserID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
	return response, nil
}
