package team

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/SenechkaP/avito-test/internal/database"
	"github.com/SenechkaP/avito-test/internal/model"
	"github.com/SenechkaP/avito-test/internal/repository"
	"github.com/jackc/pgx/v5"
)

type TeamRepository struct {
	Db      *database.DB
	Builder sq.StatementBuilderType
}

func NewTeamRepository(db *database.DB) *TeamRepository {
	return &TeamRepository{
		Db:      db,
		Builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *TeamRepository) CreateOrUpdateTeam(ctx context.Context, team *model.Team) error {
	tx, err := r.Db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var existingTeamID string
	err = tx.QueryRow(ctx, "SELECT id FROM teams WHERE team_name = $1", team.TeamName).Scan(&existingTeamID)
	if err == nil {
		return repository.ErrTeamExists
	} else if err != pgx.ErrNoRows {
		return err
	}

	for _, m := range team.Members {
		var teamIDByID *string
		if m.UserID != "" {
			err := tx.QueryRow(ctx, "SELECT team_id FROM users WHERE id = $1", m.UserID).Scan(&teamIDByID)
			if err == nil {
				if teamIDByID != nil && *teamIDByID != "" {
					return repository.ErrUserInAnotherTeam
				}
			} else if err != pgx.ErrNoRows {
				return err
			}
		}
	}

	var teamID string
	if err := tx.QueryRow(ctx, "INSERT INTO teams (team_name, created_at) VALUES ($1, now()) RETURNING id", team.TeamName).Scan(&teamID); err != nil {
		return fmt.Errorf("insert team: %w", err)
	}

	for _, m := range team.Members {
		_, err := tx.Exec(ctx, `
			INSERT INTO users (id, username, team_id, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, now(), now())
			ON CONFLICT (id) DO UPDATE SET username = EXCLUDED.username, team_id = EXCLUDED.team_id, is_active = EXCLUDED.is_active, updated_at = now()
		`, m.UserID, m.Username, teamID, m.IsActive)
		if err != nil {
			return fmt.Errorf("upsert user %s: %w", m.UserID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (r *TeamRepository) GetTeamByName(ctx context.Context, name string) (*model.Team, error) {
	sql, args, _ := r.Builder.Select("id", "team_name", "created_at").From("teams").Where(sq.Eq{"team_name": name}).ToSql()
	var t model.Team
	err := r.Db.Pool.QueryRow(ctx, sql, args...).Scan(&t.ID, &t.TeamName, &t.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	rows, err := r.Db.Pool.Query(ctx, "SELECT id, username, is_active FROM users WHERE team_id = $1 ORDER BY id", t.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	t.Members = make([]model.TeamMember, 0)
	for rows.Next() {
		var m model.TeamMember
		if err := rows.Scan(&m.UserID, &m.Username, &m.IsActive); err != nil {
			return nil, err
		}
		t.Members = append(t.Members, m)
	}
	return &t, nil
}
