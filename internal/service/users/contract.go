package users

import (
	"context"

	"github.com/SenechkaP/avito-test/internal/model"
)

type UserRepository interface {
	SetUserIsActive(ctx context.Context, userID string, isActive bool) (*model.User, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	GetPRsByReviewer(ctx context.Context, userID string) ([]model.PullRequest, error)
}
