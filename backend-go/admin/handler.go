package admin

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

func DefineAdminRoutes(admin_subrouter *mux.Router, dbtx_val dbstore.DBTX) {
	timeoutHandler := http.HandlerFunc(completeCleanDatabaseFactory(dbtx_val))
	admin_subrouter.Handle("/complete-clean", timeoutHandler).Methods(http.MethodPost)

	unverifiedDocsHandler := http.HandlerFunc(UnverifedCompleteFileSchemaListFactory(dbtx_val))
	admin_subrouter.Handle("/get-unverified-docs/{max_responses}", unverifiedDocsHandler).Methods(http.MethodGet)
}
