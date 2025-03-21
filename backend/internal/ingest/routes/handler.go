package routes

import (
	"encoding/json"
	"fmt"
<<<<<<< HEAD:backend/internal/ingest/routes/handler.go
	"kessler/internal/ingest/tasks"
=======
	"kessler/cmd/ingest/tasks"
>>>>>>> main:backend/cmd/ingest/routes/handler.go
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// @Summary		Get Version Hash
// @Description	Returns the current version hash of the application
// @Tags			system
// @Produce		plain
// @Success		200	{string}	string	"Version hash string"
// @Router			/version_hash [get]
func HandleVersionHash(w http.ResponseWriter, r *http.Request) {
	versionHash := os.Getenv("VERSION_HASH")
	w.Write([]byte(versionHash))
}

// @Summary		Add Default Ingest Task
// @Description	Creates a new default ingestion task
// @Tags			tasks
// @Accept			json
// @Produce		json
// @Param			body	body		tasks.ScraperInfoPayload	true	"Scraper information"
// @Success		200		{object}	tasks.KesslerTaskInfo
// @Failure		400		{string}	string	"Error decoding request body"
// @Failure		500		{string}	string	"Error adding task"
// @Router			/add-task/ingest [post]
func HandleDefaultIngestAddTask(w http.ResponseWriter, r *http.Request) {
	HandleIngestAddTaskGeneric[tasks.ScraperInfoPayload](w, r)
}

// @Summary		Add NYPUC Ingest Task
// @Description	Creates a new NYPUC-specific ingestion task
// @Tags			tasks
// @Accept			json
// @Produce		json
// @Param			body	body		tasks.NYPUCDocInfo	true	"NYPUC document information"
// @Success		200		{object}	tasks.KesslerTaskInfo
// @Failure		400		{string}	string	"Error decoding request body"
// @Failure		500		{string}	string	"Error adding task"
// @Router			/add-task/ingest/nypuc [post]
func HandleNYPUCIngestAddTask(w http.ResponseWriter, r *http.Request) {
	HandleIngestAddTaskGeneric[tasks.NYPUCDocInfo](w, r)
}

// @Summary		Get Task Information
// @Description	Retrieves information about a specific task by ID
// @Tags			tasks
// @Produce		json
// @Param			id	path		string	true	"Task ID"
// @Success		200	{object}	tasks.KesslerTaskInfo
// @Failure		404	{string}	string	"Error retrieving task info"
// @Failure		501	{string}	string	"Not implemented"
// @Router			/task/{id} [get]
func HandleGetTaskInfo(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func DefineGlobalRouter(global_subrouter *mux.Router) {
	global_subrouter.HandleFunc("/version_hash", HandleVersionHash).Methods(http.MethodGet)
	// global_subrouter.HandleFunc("/swaggerdata", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Printf("We are serving the thing!!!\n")
	// 	w.Header().Set("Content-Type", "application/json")
	// 	http.ServeFile(w, r, "docs/swagger.json")
	// },
	// ).Methods(http.MethodGet)
	global_subrouter.HandleFunc(
		"/add-task/process/{uuid}",
		HandleAddProcessFileTask,
	).Methods(http.MethodPost)
	global_subrouter.HandleFunc(
		"/add-task/ingest",
		HandleDefaultIngestAddTask,
	).Methods(http.MethodPost)
	global_subrouter.HandleFunc(
		"/add-task/ingest/nypuc",
		HandleNYPUCIngestAddTask,
	).Methods(http.MethodPost)
	global_subrouter.HandleFunc(
		"/add-task/ingest/openscrapers",
		HandleOpenScraperIngestAddTask,
	).Methods(http.MethodPost)
	global_subrouter.HandleFunc(
		"/task/{id}",
		HandleGetTaskInfo,
	).Methods(http.MethodGet)
}

// This is a test

// @Summary		Add Process File Task
// @Description	Creates a new task to process a file with the given UUID
// @Tags			tasks
// @Produce		json
// @Param			uuid	path		string	true	"File UUID"
// @Success		200	{object}	tasks.KesslerTaskInfo
// @Failure		400	{string}	string	"Error parsing uuid"
// @Failure		500	{string}	string	"Error adding task"
// @Router			/add-task/process/{uuid} [post]
func HandleAddProcessFileTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fileIDstr := params["uuid"]
	fileID, err := uuid.Parse(fileIDstr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing uuid: %v", err), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	payload := tasks.ProcessFilePayload{FileID: fileID}
	task, err := tasks.NewProcessFileTask(payload)
	if err != nil {
		log.Info("Error Adding Process File Task", "err", err)
		http.Error(w, fmt.Sprintf("Error adding task: %v", err), http.StatusInternalServerError)
		return
	}
	result_info, err := tasks.EnqueueTaskFromCtx(ctx, task)
	if err != nil {
		log.Info("Error Adding Process File Task", "err", err)
		http.Error(w, fmt.Sprintf("Error adding task: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result_info)
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
