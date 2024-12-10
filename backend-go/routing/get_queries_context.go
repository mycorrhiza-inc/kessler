package routing

import (
	"kessler/gen/dbstore"
	"net/http"
)

func DBQueriesFromRequest(r *http.Request) *dbstore.Queries {
	dbtx := r.Context().Value("db").(dbstore.DBTX)
	q := dbstore.New(dbtx)
	return q
}
