package repository

import "errors"

var (
	ErrTeamExists        = errors.New("team already exists")
	ErrPRExists          = errors.New("pull request already exists")
	ErrNotFound          = errors.New("not found")
	ErrPRMerged          = errors.New("pr is merged")
	ErrNotAssigned       = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate       = errors.New("no active replacement candidate in team")
	ErrUserInAnotherTeam = errors.New("user belongs to another team")
)
