package jobs

import (
	"fmt"
	"kessler/quickwit"
	"kessler/util"
	"net/http"

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
		"/index/repopulate/files",
		HandleQuickwitFileIngestFromPostgres,
	).Methods(http.MethodGet)

	// 	job_subrouter.HandleFunc(
	// 		"/new",
	// 	)

	// 	job_subrouter.HandleFunc(
	// 		"/{job_id}",
	// 	)
	// }

	// func DefineIndexJobs(parentRouter *mux.Router) {
	// 	job_subrouter := parentRouter.PathPrefix("/index").Subrouter()

	// 	job_subrouter.HandleFunc(
	// 		"/create",
	// 	)

	// job_subrouter.HandleFunc(
	// 	"/status/{job_id}",
	// 	indexJobStatusHandler,
	// )
}

func CreateConversationIndexJobHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	err := quickwit.CreateQuickwitIndexConversations()
	if err != nil {
		errorstring := fmt.Sprintf("Error creating quickwit index: %v", err)
		fmt.Println(errorstring)
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
		fmt.Println(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Organization index being created"))
}

func IndexAllDocketsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := *util.DBQueriesFromRequest(r)
	err := quickwit.IndexAllConversations(q, ctx, "")
	if err != nil {
		errorstring := fmt.Sprintf("Error ingesting dockets index: %v", err)
		fmt.Println(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Dockets being indexed"))
}

func IndexAllOrganizationsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := *util.DBQueriesFromRequest(r)
	err := quickwit.ReindexAllOrganizations(ctx, q, "")
	if err != nil {
		errorstring := fmt.Sprintf("Error ingesting orgs index: %v", err)
		fmt.Println(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Dockets being indexed"))
}

// Create a new index job
// job := jobManager.NewJob(IndexcollectionJob, nil)
// err := job.Start(indexCollection, nil)
// if err != nil {
// 	http.Error(w, "Error starting job", http.StatusInternalServerError)
// 	return
// }

// // Save the job state
// err = job.SaveState()
// if err != nil {
// 	http.Error(w, "Error saving job state", http.StatusInternalServerError)
// 	return
// }
// func indexJobStatusHandler(w http.ResponseWriter, r *http.Request) {
// 	job_id := mux.Vars(r)["job_id"]
// 	job := jobManager.GetJob(job_id)
// 	if job == nil {
// 		http.Error(w, "Job not found", http.StatusNotFound)
// 		return
// 	}

// 	status := job.GetStatus()
// 	w.Write([]byte(status))
// }
