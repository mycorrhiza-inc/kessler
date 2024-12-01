package admin

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

func DefineAdminRoutes(admin_subrouter *mux.Router, dbtx_val dbstore.DBTX) {
	timeoutHandler := http.TimeoutHandler(http.HandlerFunc(completeCleanDatabaseFactory(dbtx_val)), 10*time.Minute, "Request timed out")
	admin_subrouter.Handle("/complete-clean", timeoutHandler).Methods(http.MethodPost)

	unverifiedDocsHandler := http.TimeoutHandler(http.HandlerFunc(UnverifedCompleteFileSchemaListFactory(dbtx_val)), 10*time.Minute, "Request timed out")
	admin_subrouter.Handle("/get-unverified-docs/{max_responses}", unverifiedDocsHandler).Methods(http.MethodGet)
}
