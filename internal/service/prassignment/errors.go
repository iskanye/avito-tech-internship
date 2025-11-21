package prassignment

import "errors"

var (
	ErrNotFound   = errors.New("resource not found")
	ErrTeamExists = errors.New("team_name already exists")
	ErrPRExists   = errors.New("PR is already exists")
	ErrPRMerged   = errors.New("cannot reassign on merged PR")
)
