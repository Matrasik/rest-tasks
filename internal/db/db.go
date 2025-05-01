package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"strconv"
	"time"
)

func CreateConnPool(ctx context.Context, conn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, err
	}
	maxConns, err := strconv.Atoi(os.Getenv("MAXCONNS"))
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = int32(maxConns)
	cfg.MaxConnIdleTime = time.Minute * 30
	cfg.MaxConnLifetime = time.Hour
	pool, err := pgxpool.New(ctx, conn)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
