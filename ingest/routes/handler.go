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

func HandleVersionHash(w http.ResponseWriter, r *http.Request) {
	// Get the version hash from the environment variable
	versionHash := os.Getenv("VERSION_HASH")
	w.Write([]byte(versionHash))
}

func HandleGetTaskInfo(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// taskID := vars["id"]

	http.Error(w, "not implemented", http.StatusNotImplemented)

	// Get client from context (or create inspector in main.go)
	// client := tasks.GetClient(r.Context())

	// taskInfo, err := client.
	// if err != nil {
	// 	errorString := fmt.Sprintf("Error retrieving task info: %v", err)
	// 	log.Error(errorString)
	// 	http.Error(w, errorString, http.StatusNotFound)
	// 	return
	// }
	//
	// kessler_info := GenerateTaskInfoFromInfo(*taskInfo)
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(kessler_info)
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

func HandleDefaultIngestAddTask(w http.ResponseWriter, r *http.Request) {
	HandleIngestAddTaskGeneric[tasks.ScraperInfoPayload](w, r)
}

func HandleNYPUCIngestAddTask(w http.ResponseWriter, r *http.Request) {
	HandleIngestAddTaskGeneric[tasks.NYPUCDocInfo](w, r)
}
