package sys

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, conf *Config) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, conf.DSN)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
