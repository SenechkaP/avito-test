package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/SenechkaP/avito-test/configs"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB(ctx context.Context, cfg *configs.DbConfig) (*DB, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.GetConnectionString())
	if err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = time.Minute * 30

	var lastErr error
	for i := range cfg.AttemptsToConnect {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context canceled while connecting to db: %w", ctx.Err())
		default:
		}

		pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
		if err != nil {
			lastErr = fmt.Errorf("unable to create connection pool: %w", err)
		} else {
			if pingErr := pool.Ping(ctx); pingErr == nil {
				log.Println("Successfully connected to database")
				return &DB{Pool: pool}, nil
			} else {
				lastErr = fmt.Errorf("unable to ping database: %w", pingErr)
				pool.Close()
			}
		}

		delay := time.Duration((i + 1)) * cfg.BaseDelay
		log.Printf("DB connection attempt %d/%d failed: %v — retrying in %s", i+1, cfg.AttemptsToConnect, lastErr, delay)
		timer := time.NewTimer(delay)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return nil, fmt.Errorf("context canceled while waiting to retry db connection: %w", ctx.Err())
		}
	}

	return nil, fmt.Errorf("could not connect to database after %d attempts: %w", cfg.AttemptsToConnect, lastErr)
}

func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
