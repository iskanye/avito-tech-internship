package server

import (
	"context"

	"github.com/iskanye/avito-tech-internship/internal/api"
)

// (GET /users/getReview)
func (Server) GetUsersGetReview(
	c context.Context,
	req api.GetUsersGetReviewRequestObject,
) (api.GetUsersGetReviewResponseObject, error) {
	return nil, nil
}

// (POST /users/setIsActive)
func (Server) PostUsersSetIsActive(
	c context.Context,
	req api.PostUsersSetIsActiveRequestObject,
) (api.PostUsersSetIsActiveResponseObject, error) {
	return nil, nil
}
