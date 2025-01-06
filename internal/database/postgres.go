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
	dsn := "postgresql://neondb_owner:d7HoDfVMYLp3@ep-floral-water-a5p6auv5.us-east-2.aws.neon.tech/neondb?sslmode=require"

	var lastErr error
	for retriesLeft := maxRetries; retriesLeft > 0; retriesLeft-- {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		config, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to parse config: %w", err)
		}
		config.MaxConns = int32(cfg.MaxConn)
		config.MinConns = 1

		pool, err := pgxpool.NewWithConfig(ctx, config)
		if err == nil {
			slog.Info("Successfully connected to the database")
			return pool, nil
		}

		lastErr = err
		slog.Warn(
			"Retrying database connection",
			"attempt", maxRetries-retriesLeft+1,
			"remainingRetries", retriesLeft-1,
			"error", err,
		)
		time.Sleep(1 * time.Second) // Increase sleep duration
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
			return ctx.Err() // Return if context is done
		default:
		}

		// Attempt to ping the database
		if err := pgxPool.Ping(ctx); err == nil {
			return nil // Ping succeeded
		} else {
			lastErr = err
			slog.Info("Failed to ping database, retrying...", "attempt", i, "error", err)
			time.Sleep(500 * time.Millisecond)
		}
	}

	return fmt.Errorf("failed to ping database after retries: %w", lastErr)
}
