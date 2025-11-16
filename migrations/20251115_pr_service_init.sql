-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS teams (
    id          TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    team_name   TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users (
    id          TEXT PRIMARY KEY,
    username    TEXT NOT NULL,
    team_id     TEXT REFERENCES teams(id) ON DELETE SET NULL,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMP DEFAULT now(),
    updated_at  TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS pull_requests (
    id                TEXT PRIMARY KEY,
    pull_request_name TEXT NOT NULL,
    author_id         TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status            TEXT NOT NULL DEFAULT 'OPEN',
    created_at        TIMESTAMP DEFAULT now(),
    merged_at         TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS pr_reviewers (
    id         BIGSERIAL PRIMARY KEY,
    pr_id      TEXT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT now(),
    UNIQUE(pr_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pr_reviewers;
DROP TABLE IF EXISTS pull_requests;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams;
-- +goose StatementEnd
