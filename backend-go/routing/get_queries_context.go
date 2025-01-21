package routing

import (
	"context"
	"kessler/gen/dbstore"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DBQueriesFromRequest(r *http.Request) *dbstore.Queries {
	dbtx := r.Context().Value("db").(dbstore.DBTX)
	q := dbstore.New(dbtx)
	return q
}

//	func DBQueriesFromContext(ctx context.Context) *dbstore.Queries {
//		dbtx := ctx.Value("db").(dbstore.DBTX)
//		q := dbstore.New(dbtx)
//		return q
//	}
func DBQueriesFromContext(ctx context.Context) *dbstore.Queries {
	pool := ctx.Value("db").(*pgxpool.Pool)
	singleconn, err := pool.Acquire(ctx)
	if err != nil {
		log.Fatalf("Failed to acquire a connection from the pool: %v", err)
	}
	q := dbstore.New(singleconn)
	return q
}

func DBTXFromContext(ctx context.Context) dbstore.DBTX {
	dbtx := ctx.Value("db").(dbstore.DBTX)
	return dbtx
}
