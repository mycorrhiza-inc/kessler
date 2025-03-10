package jobs

import (
	"context"
	"fmt"
	"kessler/internal/database"
	"kessler/internal/quickwit"

	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gorilla/mux"
)

func DefineJobRoutes(parent_router *mux.Router) {
	parent_router.HandleFunc(
		"/index/create/conversations",
		CreateConversationIndexJobHandler,
	).Methods(http.MethodGet)
	parent_router.HandleFunc(
		"/index/create/organizations",
		CreateOrganizationIndexJobHandler,
	).Methods(http.MethodGet)
	parent_router.HandleFunc(
		"/index/repopulate/conversations",
		IndexAllDocketsHandler,
	).Methods(http.MethodGet)
	parent_router.HandleFunc(
		"/index/repopulate/organizations",
		IndexAllOrganizationsHandler,
	).Methods(http.MethodGet)
	parent_router.HandleFunc(
		"/index/create/files",
		CreateFileIndexJobHandlerFactory(false),
	).Methods(http.MethodGet)
	parent_router.HandleFunc(
		"/index/repopulate/files",
		HandleQuickwitFileIngestFromPostgresFactory(false),
	).Methods(http.MethodGet)
	parent_router.HandleFunc(
		"/index/create/files/test",
		CreateFileIndexJobHandlerFactory(true),
	).Methods(http.MethodGet)
	parent_router.HandleFunc(
		"/index/repopulate/files/test",
		HandleQuickwitFileIngestFromPostgresFactory(true),
	).Methods(http.MethodGet)
}

func CreateConversationIndexJobHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	err := quickwit.CreateQuickwitIndexConversations()
	if err != nil {
		errorstring := fmt.Sprintf("Error creating quickwit index: %v", err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Conversations index being created"))
}

func CreateOrganizationIndexJobHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	err := quickwit.CreateQuickwitOrganizationsIndex("") // Empty index name defaults to the production quickwit index
	if err != nil {
		errorstring := fmt.Sprintf("Error creating quickwit index: %v", err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Organization index being created"))
}

func CreateFileIndexJobHandlerFactory(isTest bool) func(http.ResponseWriter,
	*http.Request) {
	indexName := quickwit.NYPUCIndex
	if isTest {
		indexName = quickwit.TestNYPUCIndex
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// ctx := r.Context()
		err := quickwit.CreateQuickwitNYFileIndex(indexName) // Empty index name defaults to the production quickwit index
		if err != nil {
			errorstring := fmt.Sprintf("Error creating quickwit index: %v", err)
			log.Info(errorstring)
			http.Error(w, errorstring, http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Organization index being created"))
	}
}

func IndexAllDocketsHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	ctx := context.Background()
	// q := *util.DBQueriesFromContext(ctx)
	q := database.GetTx()
	err := quickwit.IndexAllConversations(*q, ctx, "")
	if err != nil {
		errorstring := fmt.Sprintf("Error ingesting dockets index: %v", err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Dockets being indexed"))
}

func IndexAllOrganizationsHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	ctx := context.Background()
	// q := *util.DBQueriesFromContext(ctx)
	q := database.GetTx()
	err := quickwit.ReindexAllOrganizations(ctx, *q, "")
	if err != nil {
		errorstring := fmt.Sprintf("Error ingesting orgs index: %v", err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Dockets being indexed"))
}
