package server

import (
	"context"

	"github.com/iskanye/avito-tech-internship/pkg/api"
)

// (POST /pullRequest/create)
func (Server) PostPullRequestCreate(
	c context.Context,
	req api.PostPullRequestCreateRequestObject,
) (api.PostPullRequestCreateResponseObject, error) {
	return nil, nil
}

// (POST /pullRequest/merge)
func (Server) PostPullRequestMerge(
	c context.Context,
	req api.PostPullRequestMergeRequestObject,
) (api.PostPullRequestMergeResponseObject, error) {
	return nil, nil
}

// (POST /pullRequest/reassign)
func (Server) PostPullRequestReassign(
	c context.Context,
	req api.PostPullRequestReassignRequestObject,
) (api.PostPullRequestReassignResponseObject, error) {
	return nil, nil
}
