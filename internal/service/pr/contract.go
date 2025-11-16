package pr

import (
	"context"

	"github.com/SenechkaP/avito-test/internal/model"
)

type PRRepository interface {
	CreatePullRequest(ctx context.Context, pr *model.PullRequest) error
	GetPullRequest(ctx context.Context, prID string) (*model.PullRequest, error)
	MergePullRequest(ctx context.Context, prID string) (*model.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID string, oldUserID string) (string, *model.PullRequest, error)
	GetAssignmentsCountByUser(ctx context.Context) ([]model.AssignmentStat, error)
	GetAssignmentsCountByPR(ctx context.Context) ([]model.AssignmentStat, error)
}
