package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	maxRetries    = 60
	retryInterval = 2 * time.Second
)

// DbMethods defines the interface for database operations.
type DbMethods interface {
	MakeNewPgxPool(ctx context.Context, maxRetries int) (*pgxpool.Pool, error)
	Ping(ctx context.Context) error
}

// DbConfig holds the configuration required to connect to the database.
type DbConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
	SSLMode  string
	pool     *pgxpool.Pool
}

// NewDbConfig creates a new instance of DbConfig.
func NewDbConfig(user, password, host string, port int, database, sslMode string) *DbConfig {
	return &DbConfig{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
		Database: database,
		SSLMode:  sslMode,
	}
}

// MakeNewPgxPool creates a new pgxpool.Pool instance with retry logic.
func (db *DbConfig) MakeNewPgxPool(ctx context.Context, maxRetries int) (*pgxpool.Pool, error) {
	dsn := "postgres://ian:newpassword@postgres:5432/witty"
	// dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", db.User, db.Password, db.Host, db.Port, db.Database)

	var lastErr error
	for retriesLeft := maxRetries; retriesLeft > 0; retriesLeft-- {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context canceled: %w", ctx.Err())
		default:
		}

		config, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to parse config: %w", err)
		}

		// Configure connection pool settings.
		config.MaxConns = 30
		config.MaxConnIdleTime = 5 * time.Minute
		config.HealthCheckPeriod = 1 * time.Minute
		config.MinConns = 1

		pool, err := pgxpool.NewWithConfig(ctx, config)
		if err == nil {
			db.pool = pool
			return pool, nil
		}

		slog.Info("Failed to connect to database, retrying", "error", err, "retriesLeft", retriesLeft-1)
		time.Sleep(500 * time.Millisecond)
		lastErr = err
	}

	return nil, fmt.Errorf("failed to connect to database after %d retries: %w", maxRetries, lastErr)
}

// Ping verifies the database connection.
func (db *DbConfig) Ping(ctx context.Context) error {
	if db.pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}
	for i := 1; i <= maxRetries; i++ {
		if err := db.pool.Ping(ctx); err != nil {
			slog.Warn(fmt.Sprintf("Ping attempt %d/%d failed: %v", i, maxRetries, err))
			if i == maxRetries {
				return fmt.Errorf("failed to ping database after %d retries: %w", maxRetries, err)

			}

			// wait before retrying
			time.Sleep(retryInterval)
			continue
		}
		slog.Info("Pinged database succesfully")
		return nil
	}
	return fmt.Errorf("unexpected error while pinging database")

}
