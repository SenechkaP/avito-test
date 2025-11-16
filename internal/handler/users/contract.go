package users

import (
	"context"

	"github.com/SenechkaP/avito-test/internal/model"
)

type UserServiceInterface interface {
	SetUserActive(ctx context.Context, userID string, isActive bool) (*model.User, error)
	GetPRsByReviewer(ctx context.Context, userID string) ([]model.PullRequest, error)
}
