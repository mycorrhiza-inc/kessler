package admin

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

// middlware to set if
func withVerifiedFileData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add the connection pool to the request context
		ctx := context.WithValue(r.Context(), "verified_search", true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func withUnverifiedFileData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add the connection pool to the request context
		ctx := context.WithValue(r.Context(), "verified_search", false)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func DefineAdminRoutes(admin_subrouter *mux.Router) {
	admin_subrouter.HandleFunc(
		"/complete-clean",
		completeCleanDatabaseHandler,
	).Methods(http.MethodPost)

	admin_subrouter.HandleFunc(
		"/get-unverified-docs/{max_responses}",
		HandleUnverifedCompleteFileSchemaList,
	).Methods(http.MethodGet)

	admin_subrouter.HandleFunc(
		"/generate-email-info/{org_uuid}",
		ExtractRelaventEmailsFromOrgUUIDHandler,
	).Methods(http.MethodGet)
	admin_subrouter.HandleFunc(
		"/file-metadata-match",
		HandleCheckDocumentMetadata,
	).Methods(http.MethodPost)

	verified := admin_subrouter.Methods(http.MethodPost).Subrouter()
	verified.HandleFunc(
		"/quickwit/ingest",
		HandleQuickwitIngestFromPostgres,
	)
	verified.Use(withUnverifiedFileData)

	unVerified := admin_subrouter.Methods(http.MethodPost).Subrouter()
	unVerified.HandleFunc(
		"/quickwit/ingest/verified_only",
		HandleQuickwitIngestFromPostgres,
	)
	unVerified.Use(withVerifiedFileData)
}
