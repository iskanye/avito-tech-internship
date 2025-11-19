package models

import "time"

const (
	PULLREQUEST_OPEN   = "OPEN"
	PULLREQUEST_MERGED = "MERGED"
)

type PullRequest struct {
	ID                string
	Name              string
	AuthorID          string
	Status            string
	AssignedReviewers []string
	CreatedAt         *time.Time
	MergedAt          *time.Time
}
