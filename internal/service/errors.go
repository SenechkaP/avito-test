package service

import "errors"

var (
	ErrTeamExists    = errors.New("team already exists")
	ErrNotFound      = errors.New("not found")
	ErrPRExists      = errors.New("pull request already exists")
	ErrPRMerged      = errors.New("cannot modify merged PR")
	ErrNotAssigned   = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate   = errors.New("no active replacement candidate in team")
	ErrInvalidParams = errors.New("invalid parameters")
)
