package server

import (
	"context"

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
func (serverAPI) PostUsersSetIsActive(
	c context.Context,
	req api.PostUsersSetIsActiveRequestObject,
) (api.PostUsersSetIsActiveResponseObject, error) {
	return nil, nil
}
