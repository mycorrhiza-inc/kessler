package util

import (
	"context"
	"kessler/gen/dbstore"
)

func DBQueriesFromContext(ctx context.Context) *dbstore.Queries {
	dbtx := ctx.Value("db").(dbstore.DBTX)
	q := dbstore.New(dbtx)
	return q
}
