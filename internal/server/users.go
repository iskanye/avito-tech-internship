package server

import (
	"context"
	"errors"

	prassignment "github.com/iskanye/avito-tech-internship/internal/service/pr_assignment"
	"github.com/iskanye/avito-tech-internship/pkg/api"
)

// (GET /users/getReview)
func (serverAPI) GetUsersGetReview(
	c context.Context,
	req api.GetUsersGetReviewRequestObject,
) (api.GetUsersGetReviewResponseObject, error) {
	return nil, nil
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
		return response, err
	}

	response := api.PostUsersSetIsActive200JSONResponse{}
	response.User = &api.User{
		UserId:   user.UserID,
		Username: user.Username,
		IsActive: user.IsActive,
	}
	return response, nil
}
