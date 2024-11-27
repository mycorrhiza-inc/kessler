package admin

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

func DefineAdminRoutes(router *mux.Router, dbtx_val dbstore.DBTX) {
	admin_subrouter := router.PathPrefix("/v2/admin").Subrouter()
	admin_subrouter.HandleFunc(
		"/complete-clean",
		completeCleanDatabaseFactory(dbtx_val)).Methods(http.MethodPost)
}
