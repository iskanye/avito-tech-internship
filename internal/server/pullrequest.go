package server

import (
	"context"

	"github.com/iskanye/avito-tech-internship/pkg/api"
)

// (POST /pullRequest/create)
func (serverAPI) PostPullRequestCreate(
	c context.Context,
	req api.PostPullRequestCreateRequestObject,
) (api.PostPullRequestCreateResponseObject, error) {
	return nil, nil
}

// (POST /pullRequest/merge)
func (serverAPI) PostPullRequestMerge(
	c context.Context,
	req api.PostPullRequestMergeRequestObject,
) (api.PostPullRequestMergeResponseObject, error) {
	return nil, nil
}

// (POST /pullRequest/reassign)
func (serverAPI) PostPullRequestReassign(
	c context.Context,
	req api.PostPullRequestReassignRequestObject,
) (api.PostPullRequestReassignResponseObject, error) {
	return nil, nil
}
