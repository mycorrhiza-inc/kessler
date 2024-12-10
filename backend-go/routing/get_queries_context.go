package routing

import (
	"context"
	"kessler/gen/dbstore"
	"net/http"
)

func DBQueriesFromRequest(r *http.Request) *dbstore.Queries {
	dbtx := r.Context().Value("db").(dbstore.DBTX)
	q := dbstore.New(dbtx)
	return q
}

func DBQueriesFromContext(ctx context.Context) *dbstore.Queries {
	dbtx := ctx.Value("db").(dbstore.DBTX)
	q := dbstore.New(dbtx)
	return q
}

func DBTXFromContext(ctx context.Context) dbstore.DBTX {
	dbtx := ctx.Value("db").(dbstore.DBTX)
	return dbtx
}
