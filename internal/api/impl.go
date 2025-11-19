package api

import "context"

type Server struct{}

func NewServer() Server {
	return Server{}
}

// (POST /team/add)
func (Server) PostTeamAdd(
	c context.Context,
	req PostTeamAddRequestObject,
) (PostTeamAddResponseObject, error) {
	return nil, nil
}

// (GET /team/get)
func (Server) GetTeamGet(
	c context.Context,
	req GetTeamGetRequestObject,
) (GetTeamGetResponseObject, error) {
	return nil, nil
}

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

// (POST /pullRequest/create)
func (Server) PostPullRequestCreate(
	c context.Context,
	req PostPullRequestCreateRequestObject,
) (PostPullRequestCreateResponseObject, error) {
	return nil, nil
}

// (POST /pullRequest/merge)
func (Server) PostPullRequestMerge(
	c context.Context,
	req PostPullRequestMergeRequestObject,
) (PostPullRequestMergeResponseObject, error) {
	return nil, nil
}

// (POST /pullRequest/reassign)
func (Server) PostPullRequestReassign(
	c context.Context,
	req PostPullRequestReassignRequestObject,
) (PostPullRequestReassignResponseObject, error) {
	return nil, nil
}
