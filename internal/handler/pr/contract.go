package pr

import (
	"context"

	"github.com/SenechkaP/avito-test/internal/model"
)

type PRServiceInterface interface {
	GetPullRequest(ctx context.Context, prID string) (*model.PullRequest, error)
	CreatePullRequest(ctx context.Context, req *model.CreatePRRequest) (*model.PullRequest, error)
	MergePullRequest(ctx context.Context, prID string) (*model.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID string, oldUserID string) (string, *model.PullRequest, error)
	GetAssignmentsCountByUser(ctx context.Context) ([]model.AssignmentStat, error)
	GetAssignmentsCountByPR(ctx context.Context) ([]model.AssignmentStat, error)
}
