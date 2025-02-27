package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"thaumaturgy/tasks"

	"github.com/charmbracelet/log"
	"github.com/gorilla/mux"
)

//	@Summary		Get Version Hash
//	@Description	Returns the current version hash of the application
//	@Tags			system
//	@Produce		plain
//	@Success		200	{string}	string	"Version hash string"
//	@Router			/version_hash [get]
func HandleVersionHash(w http.ResponseWriter, r *http.Request) {
	versionHash := os.Getenv("VERSION_HASH")
	w.Write([]byte(versionHash))
}

//	@Summary		Add Default Ingest Task
//	@Description	Creates a new default ingestion task
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			body	body		tasks.ScraperInfoPayload	true	"Scraper information"
//	@Success		200		{object}	tasks.KesslerTaskInfo
//	@Failure		400		{string}	string	"Error decoding request body"
//	@Failure		500		{string}	string	"Error adding task"
//	@Router			/add-task/ingest [post]
func HandleDefaultIngestAddTask(w http.ResponseWriter, r *http.Request) {
	HandleIngestAddTaskGeneric[tasks.ScraperInfoPayload](w, r)
}

//	@Summary		Add NYPUC Ingest Task
//	@Description	Creates a new NYPUC-specific ingestion task
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			body	body		tasks.NYPUCDocInfo	true	"NYPUC document information"
//	@Success		200		{object}	tasks.KesslerTaskInfo
//	@Failure		400		{string}	string	"Error decoding request body"
//	@Failure		500		{string}	string	"Error adding task"
//	@Router			/add-task/ingest/nypuc [post]
func HandleNYPUCIngestAddTask(w http.ResponseWriter, r *http.Request) {
	HandleIngestAddTaskGeneric[tasks.NYPUCDocInfo](w, r)
}

//	@Summary		Get Task Information
//	@Description	Retrieves information about a specific task by ID
//	@Tags			tasks
//	@Produce		json
//	@Param			id	path		string	true	"Task ID"
//	@Success		200	{object}	tasks.KesslerTaskInfo
//	@Failure		404	{string}	string	"Error retrieving task info"
//	@Failure		501	{string}	string	"Not implemented"
//	@Router			/task/{id} [get]
func HandleGetTaskInfo(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func DefineGlobalRouter(global_subrouter *mux.Router) {
	global_subrouter.HandleFunc("/version_hash", HandleVersionHash).Methods(http.MethodGet)
	global_subrouter.HandleFunc(
		"/add-task/ingest",
		HandleDefaultIngestAddTask,
	).Methods(http.MethodPost)
	global_subrouter.HandleFunc(
		"/add-task/ingest/nypuc",
		HandleNYPUCIngestAddTask,
	).Methods(http.MethodPost)
	global_subrouter.HandleFunc(
		"/task/{id}",
		HandleGetTaskInfo,
	).Methods(http.MethodGet)
}

func HandleIngestAddTaskGeneric[T tasks.CastableIntoScraperInfo](w http.ResponseWriter, r *http.Request) {
	var scraper_info T
	if err := json.NewDecoder(r.Body).Decode(&scraper_info); err != nil {
		log.Info("User Gave Bad Request", "err", err)
		errorString := fmt.Sprintf("Error decoding request body: %v", err)
		http.Error(w, errorString, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	kessler_info, err := tasks.AddScraperTaskCastable(ctx, scraper_info)
	if err != nil {
		log.Error("Encountered Error Adding Task", "err", err)
		http.Error(w, fmt.Sprintf("Error adding task: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kessler_info)
}
