package pgclient

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrFailedConnectionToDatabase = fmt.Errorf("failed connection to database")
)

type Config struct {
	Username string
	Password string
	Host string
	Port string
	Database string
}

func NewClient(ctx context.Context, connAttempts int, config Config) (*pgxpool.Pool, error) {
	slog.Info("pool creation started")
	
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.Username, config.Password, config.Host, config.Port, config.Database)
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		connAttempts = 0
	}

	for i := 0; i < connAttempts; i++ {
		ctx, cancel := context.WithTimeout(ctx, 5 * time.Second)
		defer cancel()
		
		slog.Info(fmt.Sprintf("attempt %d to connect to the database", i))
		err = pool.Ping(ctx)
		if err != nil {
			slog.Warn("failed to connect to the database")
			time.Sleep(5 * time.Second)
			continue
		}

		slog.Info("successful connection to the database")
		slog.Info("pool creation successful")
		return pool, nil
	}

	slog.Error("pool creation failed")
	return nil, ErrFailedConnectionToDatabase
}
