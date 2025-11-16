package users

import (
	"context"
	"errors"

	"github.com/SenechkaP/avito-test/internal/model"
	"github.com/SenechkaP/avito-test/internal/repository"
	"github.com/SenechkaP/avito-test/internal/service"
)

type UserService struct {
	Repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) SetUserActive(ctx context.Context, userID string, isActive bool) (*model.User, error) {
	u, err := s.Repo.SetUserIsActive(ctx, userID, isActive)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, service.ErrNotFound
		}
		return nil, err
	}
	return u, nil
}

func (s *UserService) GetPRsByReviewer(ctx context.Context, userID string) ([]model.PullRequest, error) {
	if userID == "" {
		return nil, service.ErrInvalidParams
	}
	prs, err := s.Repo.GetPRsByReviewer(ctx, userID)
	if err != nil {
		return nil, err
	}
	return prs, nil
}
