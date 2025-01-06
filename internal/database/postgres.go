package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DbConfig holds the database configuration.
type DbConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
	MaxConn  int
}

// DBMethods defines the interface for database operations.
type DBMethods interface {
	NewPgxPool(ctx context.Context, maxRetries int) (*pgxpool.Pool, error)
	Ping(ctx context.Context, pgxPool *pgxpool.Pool, maxRetries int) error
}

// NewPgxPool creates a new pgx connection pool with retries.
func (cfg *DbConfig) NewPgxPool(ctx context.Context, maxRetries int) (*pgxpool.Pool, error) {
	// Construct DSN
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=require",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)

	var lastErr error
	for retriesLeft := maxRetries; retriesLeft > 0; retriesLeft-- {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context canceled: %w", ctx.Err())
		default:
			// Parse configuration
			config, err := pgxpool.ParseConfig(dsn)
			if err != nil {
				return nil, fmt.Errorf("failed to parse config: %w", err)
			}
			config.MaxConns = int32(cfg.MaxConn)
			config.MinConns = 1

			// Attempt to create a connection pool
			pool, err := pgxpool.NewWithConfig(ctx, config)
			if err == nil {
				return pool, nil
			}

			lastErr = err
			slog.Info("Failed to connect to database, retrying...", "error", err, "retriesLeft", retriesLeft-1)
			time.Sleep(500 * time.Millisecond)
		}
	}

	return nil, fmt.Errorf("failed to connect to database after retries: %w", lastErr)
}

// Ping checks the health of the database connection with retries.
func (cfg *DbConfig) Ping(ctx context.Context, pgxPool *pgxpool.Pool, maxRetries int) error {
	if cfg == nil {
		return fmt.Errorf("db config is nil")
	}
	if pgxPool == nil {
		return fmt.Errorf("pgxPool is nil")
	}

	var lastErr error
	for i := 1; i <= maxRetries; i++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled during ping: %w", ctx.Err())
		default:
			// Attempt to ping the database
			if err := pgxPool.Ping(ctx); err == nil {
				return nil // Ping succeeded
			} else {
				lastErr = err
				slog.Info("Failed to ping database, retrying...", "attempt", i, "error", err)
				time.Sleep(500 * time.Millisecond)
			}
		}
	}

	return fmt.Errorf("failed to ping database after retries: %w", lastErr)
}
