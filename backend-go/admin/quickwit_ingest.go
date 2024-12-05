package admin

import (
	"kessler/gen/dbstore"
	"kessler/quickwit"
	"net/http"
)

func HandleQuickwitIngestFromPostgresFactory(dbtx_val dbstore.DBTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := *dbstore.New(dbtx_val)
		indexName := "idk what the right index name is"
		quickwit.ClearIndex(indexName)
	}
}
