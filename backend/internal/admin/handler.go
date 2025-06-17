package admin

import (
	"context"
	"kessler/internal/dbstore"
	"net/http"

	"github.com/gorilla/mux"
)

var tracer = otel.Tracer("admin-handler")

type AdminHandler struct {
	db dbstore.DBTX
}

func NewAdminHandler(db dbstore.DBTX) *AdminHandler {
	return &AdminHandler{
		db,
	}
}

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

func DefineAdminRoutes(admin_subrouter *mux.Router, db dbstore.DBTX) {
	handler := NewAdminHandler(db)
	admin_subrouter.HandleFunc(
		"/complete-clean",
		handler.completeCleanDatabaseHandler,
	).Methods(http.MethodPost)

	admin_subrouter.HandleFunc(
		"/get-unverified-docs/{max_responses}",
		handler.HandleUnverifedCompleteFileSchemaList,
	).Methods(http.MethodGet)

	// admin_subrouter.HandleFunc(
	// 	"/generate-email-info/{org_uuid}",
	// 	ExtractRelaventEmailsFromOrgUUIDHandler,
	// ).Methods(http.MethodGet)
	admin_subrouter.HandleFunc(
		"/file-metadata-match",
		handler.HandleCheckDocumentMetadata,
	).Methods(http.MethodPost)
}
