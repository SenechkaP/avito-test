package team

import (
	"context"

	"github.com/SenechkaP/avito-test/internal/model"
)

type TeamServiceInterface interface {
	CreateOrUpdateTeam(ctx context.Context, req *model.CreateTeamRequest) error
	GetTeamByName(ctx context.Context, name string) (*model.Team, error)
}
