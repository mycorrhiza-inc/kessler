package util

import (
	"context"
	"kessler/gen/dbstore"

	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DBQueriesFromContext(ctx context.Context) *dbstore.Queries {
	pool := ctx.Value("db").(*pgxpool.Pool)
	singleconn, err := pool.Acquire(ctx)
	if err != nil {
		log.Fatalf("Failed to acquire a connection from the pool: %v", err)
	}
	q := dbstore.New(singleconn)
	return q
}
