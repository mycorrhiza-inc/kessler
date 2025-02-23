package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"thaumaturgy/tasks"

	"github.com/charmbracelet/log"
	"github.com/gorilla/mux"
)

func DefineGlobalRouter(global_subrouter *mux.Router) {
	global_subrouter.HandleFunc(
		"/add-task/ingest",
		HandleGenericIngestAddTask,
	).Methods(http.MethodPost)
	global_subrouter.HandleFunc(
		"/task/{id}",
		HandleGetTaskInfo,
	).Methods(http.MethodGet)
}

func HandleGetTaskInfo(w http.ResponseWriter, r *http.Request) {
}

func HandleGenericIngestAddTask(w http.ResponseWriter, r *http.Request) {
	var scraper_info tasks.ScraperInfoPayload
	body_bytes, err := io.ReadAll(r.Body)
	if err != nil {
		error_string := fmt.Sprintf("Error reading request body: %v", err)
		log.Error(error_string)
		http.Error(w, error_string, http.StatusBadRequest)
	}
	err = json.Unmarshal(body_bytes, &scraper_info)
	if err != nil {
		error_string := fmt.Sprintf("Response not properly formatted: %v", err)
		log.Error(error_string)
		http.Error(w, error_string, http.StatusBadRequest)
	}
	task, err := tasks.NewAddFileScraperTask(scraper_info)
	if err != nil {
		error_string := fmt.Sprintf("Error creating task: %v", err)
		log.Error(error_string)
		http.Error(w, error_string, http.StatusBadRequest)
	}
}
