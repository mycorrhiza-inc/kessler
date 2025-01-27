package jobs

import (
	"fmt"
	"kessler/quickwit"
	"net/http"

	"github.com/gorilla/mux"
)

func DefineJobRoutes(parentRouter *mux.Router) {
	job_subrouter := parentRouter.PathPrefix("/jobs").Subrouter()

	job_subrouter.HandleFunc(
		"/create",
		CreateConversationIndexJobHandler,
	).Methods(http.MethodGet)
	fmt.Println("index subroute added")

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
		http.Error(w, "Error creating conversations index", http.StatusInternalServerError)
		return
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

	w.Write([]byte("Conversations index being created"))
}

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
