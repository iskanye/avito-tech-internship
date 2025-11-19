package api

import "context"

// (GET /users/getReview)
func (Server) GetUsersGetReview(
	c context.Context,
	req GetUsersGetReviewRequestObject,
) (GetUsersGetReviewResponseObject, error) {
	return nil, nil
}

// (POST /users/setIsActive)
func (Server) PostUsersSetIsActive(
	c context.Context,
	req PostUsersSetIsActiveRequestObject,
) (PostUsersSetIsActiveResponseObject, error) {
	return nil, nil
}
