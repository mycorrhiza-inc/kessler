package jobs

import (
	"context"
	"fmt"
	"kessler/internal/dbstore"
	"kessler/internal/quickwit"
	"kessler/pkg/database"

	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gorilla/mux"
)

type JobsHandler struct {
	db dbstore.DBTX
}

func NewJobsHandler(db dbstore.DBTX) *JobsHandler {
	return &JobsHandler{
		db,
	}
}

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

func (h *JobsHandler) CreateConversationIndexJobHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "jobs:CreateConversationIndexJobHandler")
	defer span.End()

	err := quickwit.CreateQuickwitIndexConversations()
	if err != nil {
		errorstring := fmt.Sprintf("Error creating quickwit index: %v", err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Conversations index being created"))
}

func (h *JobsHandler) CreateOrganizationIndexJobHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "jobs:CreateOrganizationIndexJobHandler")
	defer span.End()

	err := quickwit.CreateQuickwitOrganizationsIndex("") // Empty index name defaults to the production quickwit index
	if err != nil {
		errorstring := fmt.Sprintf("Error creating quickwit index: %v", err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Organization index being created"))
}

func (h *JobsHandler) CreateFileIndexJobHandlerFactory(isTest bool) func(http.ResponseWriter,
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

func (h *JobsHandler) IndexAllDocketsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "jobs:IndexAllDocketsHandler")
	defer span.End()

	q := database.GetQueries(h.db)
	err := quickwit.IndexAllConversations(*q, ctx, "")
	if err != nil {
		errorstring := fmt.Sprintf("Error ingesting dockets index: %v", err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Dockets being indexed"))
}

func (h *JobsHandler) IndexAllOrganizationsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "jobs:IndexAllOrganizationsHandler")
	defer span.End()

	q := database.GetQueries(h.db)
	err := quickwit.ReindexAllOrganizations(ctx, *q, "")
	if err != nil {
		errorstring := fmt.Sprintf("Error ingesting orgs index: %v", err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Dockets being indexed"))
}
