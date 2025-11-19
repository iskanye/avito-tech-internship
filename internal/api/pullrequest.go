package api

import "context"

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
