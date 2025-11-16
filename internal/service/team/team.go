package team

import (
	"context"
	"errors"
	"fmt"

	"github.com/SenechkaP/avito-test/internal/model"
	"github.com/SenechkaP/avito-test/internal/repository"
	"github.com/SenechkaP/avito-test/internal/service"
)

type TeamService struct {
	Repo TeamRepository
}

func NewTeamService(repo TeamRepository) *TeamService {
	return &TeamService{Repo: repo}
}

func (s *TeamService) CreateOrUpdateTeam(ctx context.Context, req *model.CreateTeamRequest) error {
	if req.TeamName == "" {
		return service.ErrInvalidParams
	}
	team := &model.Team{
		TeamName: req.TeamName,
	}
	team.Members = append(team.Members, req.Members...)
	if err := s.Repo.CreateOrUpdateTeam(ctx, team); err != nil {
		if errors.Is(err, repository.ErrTeamExists) {
			return service.ErrTeamExists
		}
		if errors.Is(err, repository.ErrUserInAnotherTeam) {
			return service.ErrInvalidParams
		}
		return fmt.Errorf("create team: %w", err)
	}
	return nil
}

func (s *TeamService) GetTeamByName(ctx context.Context, name string) (*model.Team, error) {
	if name == "" {
		return nil, service.ErrInvalidParams
	}
	t, err := s.Repo.GetTeamByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, service.ErrNotFound
	}
	return t, nil
}
