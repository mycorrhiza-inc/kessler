package util

import (
	"context"
	"kessler/gen/dbstore"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DBTXFromContext(ctx context.Context) dbstore.DBTX {
	pool := ctx.Value("db").(*pgxpool.Pool)
	// connection, err := pool.Acquire(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	return pool
}

func DBQueriesFromContext(ctx context.Context) *dbstore.Queries {
	dbtx := DBTXFromContext(ctx)
	q := dbstore.New(dbtx)
	return q
}

func DBQueriesFromRequest(r *http.Request) *dbstore.Queries {
	ctx := r.Context()
	q := DBQueriesFromContext(ctx)
	return q
}
