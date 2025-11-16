package team

import (
	"context"

	"github.com/SenechkaP/avito-test/internal/model"
)

type TeamRepository interface {
	CreateOrUpdateTeam(ctx context.Context, team *model.Team) error
	GetTeamByName(ctx context.Context, name string) (*model.Team, error)
}
