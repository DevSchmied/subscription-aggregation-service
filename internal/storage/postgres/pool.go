package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/DevSchmied/subscription-aggregation-service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool initializes PostgreSQL connection pool.
func NewPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {

	// Build database connection string
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	// Parse DSN into pool config
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse pg config: %w", err)
	}

	// Configure pool limits
	poolCfg.MaxConns = 10
	poolCfg.MinConns = 2
	poolCfg.MaxConnLifetime = 30 * time.Minute
	poolCfg.MaxConnIdleTime = 5 * time.Minute

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create pg pool: %w", err)
	}

	// Ping database with timeout
	ctxPing, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := pool.Ping(ctxPing); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping pg failed: %w", err)
	}

	// Return ready pool
	return pool, nil
}
