package pr

import (
	"context"
	"errors"
	"fmt"

	"github.com/SenechkaP/avito-test/internal/model"
	"github.com/SenechkaP/avito-test/internal/repository"
	"github.com/SenechkaP/avito-test/internal/service"
)

type PRService struct {
	Repo PRRepository
}

func NewPRService(repo PRRepository) *PRService {
	return &PRService{Repo: repo}
}

func (s *PRService) CreatePullRequest(ctx context.Context, req *model.CreatePRRequest) (*model.PullRequest, error) {
	if req.PullRequestID == "" || req.PullRequestName == "" || req.AuthorID == "" {
		return nil, service.ErrInvalidParams
	}
	pr := &model.PullRequest{
		ID:              req.PullRequestID,
		PullRequestName: req.PullRequestName,
		AuthorID:        req.AuthorID,
	}
	if err := s.Repo.CreatePullRequest(ctx, pr); err != nil {
		if errors.Is(err, repository.ErrPRExists) {
			return nil, service.ErrPRExists
		}
		if errors.Is(err, repository.ErrNotFound) {
			return nil, service.ErrNotFound
		}
		return nil, fmt.Errorf("create pr: %w", err)
	}
	return pr, nil
}

func (s *PRService) GetPullRequest(ctx context.Context, prID string) (*model.PullRequest, error) {
	if prID == "" {
		return nil, service.ErrInvalidParams
	}
	pr, err := s.Repo.GetPullRequest(ctx, prID)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, service.ErrNotFound
	}
	return pr, nil
}

func (s *PRService) MergePullRequest(ctx context.Context, prID string) (*model.PullRequest, error) {
	pr, err := s.Repo.MergePullRequest(ctx, prID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, service.ErrNotFound
		}
		return nil, err
	}
	return pr, nil
}

func (s *PRService) ReassignReviewer(ctx context.Context, prID string, oldUserID string) (string, *model.PullRequest, error) {
	newUID, pr, err := s.Repo.ReassignReviewer(ctx, prID, oldUserID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return "", nil, service.ErrNotFound
		case errors.Is(err, repository.ErrPRMerged):
			return "", nil, service.ErrPRMerged
		case errors.Is(err, repository.ErrNotAssigned):
			return "", nil, service.ErrNotAssigned
		case errors.Is(err, repository.ErrNoCandidate):
			return "", nil, service.ErrNoCandidate
		default:
			return "", nil, err
		}
	}
	return newUID, pr, nil
}

func (s *PRService) GetAssignmentsCountByUser(ctx context.Context) ([]model.AssignmentStat, error) {
	stats, err := s.Repo.GetAssignmentsCountByUser(ctx)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (s *PRService) GetAssignmentsCountByPR(ctx context.Context) ([]model.AssignmentStat, error) {
	stats, err := s.Repo.GetAssignmentsCountByPR(ctx)
	if err != nil {
		return nil, err
	}
	return stats, nil
}
