package users

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/SenechkaP/avito-test/internal/database"
	"github.com/SenechkaP/avito-test/internal/model"
	"github.com/SenechkaP/avito-test/internal/repository"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	Db      *database.DB
	Builder sq.StatementBuilderType
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{
		Db:      db,
		Builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *UserRepository) SetUserIsActive(ctx context.Context, userID string, isActive bool) (*model.User, error) {
	query := r.Builder.Update("users").Set("is_active", isActive).Set("updated_at", sq.Expr("now()")).Where(sq.Eq{"id": userID}).Suffix("RETURNING id, username, team_id, is_active, created_at, updated_at")
	sql, args, _ := query.ToSql()

	var u model.User
	err := r.Db.Pool.QueryRow(ctx, sql, args...).Scan(&u.ID, &u.Username, &u.TeamID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	q := r.Builder.Select("id", "username", "team_id", "is_active", "created_at", "updated_at").
		From("users").
		Where(sq.Eq{"id": userID})
	sql, args, _ := q.ToSql()

	var u model.User
	err := r.Db.Pool.QueryRow(ctx, sql, args...).Scan(&u.ID, &u.Username, &u.TeamID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetPRsByReviewer(ctx context.Context, userID string) ([]model.PullRequest, error) {
	rows, err := r.Db.Pool.Query(ctx, `
		SELECT p.id, p.pull_request_name, p.author_id, p.status, p.created_at, p.merged_at
		FROM pull_requests p
		JOIN pr_reviewers pr ON p.id = pr.pr_id
		WHERE pr.user_id = $1
		ORDER BY p.id
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prs := make([]model.PullRequest, 0)
	for rows.Next() {
		var pr model.PullRequest
		var mergedAt *time.Time
		if err := rows.Scan(&pr.ID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &mergedAt); err != nil {
			return nil, err
		}
		pr.MergedAt = mergedAt
		prs = append(prs, pr)
	}
	return prs, nil
}
