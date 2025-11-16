package pr

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/SenechkaP/avito-test/internal/database"
	"github.com/SenechkaP/avito-test/internal/model"
	"github.com/SenechkaP/avito-test/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PRRepository struct {
	Db      *database.DB
	Builder sq.StatementBuilderType
}

func NewPRRepository(db *database.DB) *PRRepository {
	return &PRRepository{
		Db:      db,
		Builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *PRRepository) CreatePullRequest(ctx context.Context, pr *model.PullRequest) error {
	tx, err := r.Db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var exists int
	if err := tx.QueryRow(ctx, "SELECT 1 FROM pull_requests WHERE id = $1", pr.ID).Scan(&exists); err == nil {
		return repository.ErrPRExists
	} else if err != pgx.ErrNoRows {
		return err
	}

	var authorTeamID string
	err = tx.QueryRow(ctx, "SELECT team_id FROM users WHERE id = $1", pr.AuthorID).Scan(&authorTeamID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return repository.ErrNotFound
		}
		return err
	}

	rows, err := tx.Query(ctx, `
        SELECT id FROM users
        WHERE team_id = $1 AND is_active = true AND id <> $2
        ORDER BY random()
        LIMIT 2
    `, authorTeamID, pr.AuthorID)
	if err != nil {
		return err
	}
	defer rows.Close()

	reviewers := make([]string, 0, 2)
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return err
		}
		reviewers = append(reviewers, uid)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO pull_requests (id, pull_request_name, author_id, status, created_at)
		VALUES ($1, $2, $3, $4, now())
	`, pr.ID, pr.PullRequestName, pr.AuthorID, model.PRStatusOpen)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return repository.ErrPRExists
		}
		return fmt.Errorf("insert pr: %w", err)
	}

	for _, uid := range reviewers {
		_, err := tx.Exec(ctx, `INSERT INTO pr_reviewers (pr_id, user_id, assigned_at) VALUES ($1, $2, now())`, pr.ID, uid)
		if err != nil {
			return fmt.Errorf("insert pr_reviewer: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	pr.AssignedReviewers = reviewers
	pr.Status = model.PRStatusOpen
	pr.CreatedAt = time.Now()
	return nil
}

func (r *PRRepository) GetPullRequest(ctx context.Context, prID string) (*model.PullRequest, error) {
	q := r.Builder.Select("id", "pull_request_name", "author_id", "status", "created_at", "merged_at").
		From("pull_requests").
		Where(sq.Eq{"id": prID})
	sql, args, _ := q.ToSql()

	var pr model.PullRequest
	var mergedAt *time.Time
	err := r.Db.Pool.QueryRow(ctx, sql, args...).Scan(&pr.ID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &mergedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	pr.MergedAt = mergedAt

	rows, err := r.Db.Pool.Query(ctx, "SELECT user_id FROM pr_reviewers WHERE pr_id = $1 ORDER BY id", pr.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pr.AssignedReviewers = []string{}
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		pr.AssignedReviewers = append(pr.AssignedReviewers, uid)
	}
	return &pr, nil
}

func (r *PRRepository) MergePullRequest(ctx context.Context, prID string) (*model.PullRequest, error) {
	tx, err := r.Db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var pr model.PullRequest
	var mergedAt *time.Time
	err = tx.QueryRow(ctx, `
		UPDATE pull_requests
		SET status = $2, merged_at = COALESCE(merged_at, now())
		WHERE id = $1
		RETURNING id, pull_request_name, author_id, status, created_at, merged_at
	`, prID, model.PRStatusMerged).Scan(&pr.ID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &mergedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	pr.MergedAt = mergedAt

	rows, err := tx.Query(ctx, "SELECT user_id FROM pr_reviewers WHERE pr_id = $1 ORDER BY id", pr.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pr.AssignedReviewers = []string{}
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		pr.AssignedReviewers = append(pr.AssignedReviewers, uid)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &pr, nil
}

func (r *PRRepository) ReassignReviewer(ctx context.Context, prID string, oldUserID string) (string, *model.PullRequest, error) {
	tx, err := r.Db.Pool.Begin(ctx)
	if err != nil {
		return "", nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var status string
	var authorID string
	if err := tx.QueryRow(ctx, "SELECT status, author_id FROM pull_requests WHERE id = $1", prID).Scan(&status, &authorID); err != nil {
		if err == pgx.ErrNoRows {
			return "", nil, repository.ErrNotFound
		}
		return "", nil, err
	}
	if status == model.PRStatusMerged {
		return "", nil, repository.ErrPRMerged
	}

	var tmp int
	if err := tx.QueryRow(ctx, "SELECT 1 FROM pr_reviewers WHERE pr_id = $1 AND user_id = $2", prID, oldUserID).Scan(&tmp); err != nil {
		if err == pgx.ErrNoRows {
			return "", nil, repository.ErrNotAssigned
		}
		return "", nil, err
	}

	var teamID string
	if err := tx.QueryRow(ctx, "SELECT team_id FROM users WHERE id = $1", oldUserID).Scan(&teamID); err != nil {
		if err == pgx.ErrNoRows {
			return "", nil, repository.ErrNotFound
		}
		return "", nil, err
	}

	row := tx.QueryRow(ctx, `
		SELECT id FROM users
		WHERE team_id = $1 AND is_active = true AND id NOT IN (
			SELECT user_id FROM pr_reviewers WHERE pr_id = $2
		) AND id <> $3
		ORDER BY random()
		LIMIT 1
	`, teamID, prID, authorID)

	var newUserID string
	if err := row.Scan(&newUserID); err != nil {
		if err == pgx.ErrNoRows {
			return "", nil, repository.ErrNoCandidate
		}
		return "", nil, err
	}

	_, err = tx.Exec(ctx, "DELETE FROM pr_reviewers WHERE pr_id = $1 AND user_id = $2", prID, oldUserID)
	if err != nil {
		return "", nil, err
	}
	_, err = tx.Exec(ctx, "INSERT INTO pr_reviewers (pr_id, user_id, assigned_at) VALUES ($1, $2, now())", prID, newUserID)
	if err != nil {
		return "", nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", nil, err
	}

	pr, err := r.GetPullRequest(ctx, prID)
	if err != nil {
		return "", nil, err
	}

	return newUserID, pr, nil
}

func (r *PRRepository) GetAssignmentsCountByUser(ctx context.Context) ([]model.AssignmentStat, error) {
	rows, err := r.Db.Pool.Query(ctx, `
		SELECT user_id, COUNT(*) AS cnt
		FROM pr_reviewers
		GROUP BY user_id
		ORDER BY cnt DESC, user_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.AssignmentStat
	for rows.Next() {
		var uid string
		var cnt int64
		if err := rows.Scan(&uid, &cnt); err != nil {
			return nil, err
		}
		out = append(out, model.AssignmentStat{Key: uid, Count: int(cnt)})
	}
	return out, nil
}

func (r *PRRepository) GetAssignmentsCountByPR(ctx context.Context) ([]model.AssignmentStat, error) {
	rows, err := r.Db.Pool.Query(ctx, `
		SELECT pr_id, COUNT(*) AS cnt
		FROM pr_reviewers
		GROUP BY pr_id
		ORDER BY cnt DESC, pr_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.AssignmentStat
	for rows.Next() {
		var pid string
		var cnt int64
		if err := rows.Scan(&pid, &cnt); err != nil {
			return nil, err
		}
		out = append(out, model.AssignmentStat{Key: pid, Count: int(cnt)})
	}
	return out, nil
}
