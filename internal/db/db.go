package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateConnPool(ctx context.Context, conn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, conn)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
