package server

import (
	"context"
	"errors"
	"time"

	"github.com/iskanye/avito-tech-internship/internal/models"
	"github.com/iskanye/avito-tech-internship/internal/service/prassignment"
	"github.com/iskanye/avito-tech-internship/pkg/api"
)

// (POST /pullRequest/create)
func (s *serverAPI) PostPullRequestCreate(
	c context.Context,
	req api.PostPullRequestCreateRequestObject,
) (api.PostPullRequestCreateResponseObject, error) {
	pullRequest := models.PullRequest{
		ID:        req.Body.PullRequestId,
		Name:      req.Body.PullRequestName,
		AuthorID:  req.Body.AuthorId,
		Status:    models.PULLREQUEST_OPEN,
		CreatedAt: time.Now(),
	}

	pullRequest, err := s.assign.CreatePullRequest(c, pullRequest)
	if errors.Is(err, prassignment.ErrNotFound) {
		response := api.PostPullRequestCreate404JSONResponse{}
		response.Error.Code = api.NOTFOUND
		response.Error.Message = err.Error()
		return response, nil
	} else if errors.Is(err, prassignment.ErrPRExists) {
		response := api.PostPullRequestCreate409JSONResponse{}
		response.Error.Code = api.PREXISTS
		response.Error.Message = err.Error()
		return response, nil
	} else if err != nil {
		return nil, err
	}

	response := api.PostPullRequestCreate201JSONResponse{
		Pr: &api.PullRequest{
			PullRequestId:     pullRequest.ID,
			PullRequestName:   pullRequest.Name,
			AuthorId:          pullRequest.AuthorID,
			Status:            api.PullRequestStatus(pullRequest.Status),
			AssignedReviewers: pullRequest.AssignedReviewers,
			CreatedAt:         &pullRequest.CreatedAt,
		},
	}

	return response, nil
}

// (POST /pullRequest/merge)
func (s *serverAPI) PostPullRequestMerge(
	c context.Context,
	req api.PostPullRequestMergeRequestObject,
) (api.PostPullRequestMergeResponseObject, error) {
	pullRequest, err := s.assign.MergePullRequest(c, req.Body.PullRequestId)
	if errors.Is(err, prassignment.ErrNotFound) {
		response := api.PostPullRequestMerge404JSONResponse{}
		response.Error.Code = api.NOTFOUND
		response.Error.Message = err.Error()
		return response, nil
	} else if err != nil {
		return nil, err
	}

	response := api.PostPullRequestMerge200JSONResponse{
		Pr: &api.PullRequest{
			PullRequestId:     pullRequest.ID,
			PullRequestName:   pullRequest.Name,
			AuthorId:          pullRequest.AuthorID,
			Status:            api.PullRequestStatus(pullRequest.Status),
			AssignedReviewers: pullRequest.AssignedReviewers,
			CreatedAt:         &pullRequest.CreatedAt,
			MergedAt:          &pullRequest.MergedAt,
		},
	}

	return response, nil
}

// (POST /pullRequest/reassign)
func (s *serverAPI) PostPullRequestReassign(
	c context.Context,
	req api.PostPullRequestReassignRequestObject,
) (api.PostPullRequestReassignResponseObject, error) {
	pullRequest, replacedBy, err := s.assign.ReassignPullRequest(c, req.Body.PullRequestId, req.Body.OldUserId)
	if errors.Is(err, prassignment.ErrNotFound) {
		response := api.PostPullRequestReassign404JSONResponse{}
		response.Error.Code = api.NOTFOUND
		response.Error.Message = err.Error()
		return response, nil
	} else if errors.Is(err, prassignment.ErrPRMerged) {
		response := api.PostPullRequestReassign409JSONResponse{}
		response.Error.Code = api.PRMERGED
		response.Error.Message = err.Error()
		return response, nil
	} else if err != nil {
		return nil, err
	}

	response := api.PostPullRequestReassign200JSONResponse{
		Pr: api.PullRequest{
			PullRequestId:     pullRequest.ID,
			PullRequestName:   pullRequest.Name,
			AuthorId:          pullRequest.AuthorID,
			Status:            api.PullRequestStatus(pullRequest.Status),
			AssignedReviewers: pullRequest.AssignedReviewers,
			CreatedAt:         &pullRequest.CreatedAt,
			MergedAt:          &pullRequest.MergedAt,
		},
		ReplacedBy: replacedBy,
	}

	return response, nil
}
