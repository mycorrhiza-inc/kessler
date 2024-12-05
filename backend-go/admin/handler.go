package admin

import (
	"kessler/gen/dbstore"
	"net/http"

	"github.com/gorilla/mux"
)

func DefineAdminRoutes(admin_subrouter *mux.Router, dbtx_val dbstore.DBTX) {
	cleanDatabaseHandler := http.HandlerFunc(completeCleanDatabaseFactory(dbtx_val))
	admin_subrouter.Handle("/complete-clean", cleanDatabaseHandler).Methods(http.MethodPost)

	unverifiedDocsHandler := http.HandlerFunc(UnverifedCompleteFileSchemaListFactory(dbtx_val))
	admin_subrouter.Handle("/get-unverified-docs/{max_responses}", unverifiedDocsHandler).Methods(http.MethodGet)

	quickwitIngestHandler := http.HandlerFunc(HandleQuickwitIngestFromPostgresFactory(dbtx_val))
	admin_subrouter.Handle("/quickwit/ingest", quickwitIngestHandler).Methods(http.MethodPost)
}
