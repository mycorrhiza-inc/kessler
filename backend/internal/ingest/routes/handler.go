package routes

import (
	"encoding/json"
	"fmt"
	"kessler/internal/ingest/tasks"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
)

// @Summary	Get Version Hash
// @Description	Returns the current version hash of the application
// @Tags		system
// @Produce	plain
// @Success	200	{string}	string	"Version hash string"
// @Router	/version_hash [get]
func HandleVersionHash(w http.ResponseWriter, r *http.Request) {
	versionHash := os.Getenv("VERSION_HASH")
	w.Write([]byte(versionHash))
}

// @Summary	Add Default Ingest Task
// @Description	Creates a new default ingestion task. This should use the default openscrapers schema.
// @Tags		tasks
// @Accept	json
// @Produce	json
// @Param	body	body		tasks.FilingInfoPayload	true	"Filing + Case information"
// @Success	200	{object}	tasks.KesslerTaskInfo
// @Failure	400	{string}	string	"Error decoding request body"
// @Failure	500	{string}	string	"Error adding task"
// @Router	/add-task/ingest [post]
func HandleDefaultIngestAddTask(w http.ResponseWriter, r *http.Request) {
	HandleIngestAddTaskGeneric[tasks.FilingInfoPayload](w, r)
}

// @Summary	Add NYPUC Ingest Task
// @Description	Creates a new NYPUC-specific ingestion task
// @Tags		tasks
// @Accept	json
// @Produce	json
// @Param	body	body		tasks.NYPUCDocInfo		true	"NYPUC document information"
// @Success	200	{object}	tasks.KesslerTaskInfo
// @Failure	400	{string}	string	"Error decoding request body"
// @Failure	500	{string}	string	"Error adding task"
// @Router	/add-task/ingest/nypuc [post]
func HandleNYPUCIngestAddTask(w http.ResponseWriter, r *http.Request) {
	HandleIngestAddTaskGeneric[tasks.NYPUCDocInfo](w, r)
}

// @Summary	Add Case Ingest Task
// @Description	Creates a new Case-specific ingestion task
// @Tags		tasks
// @Accept	json
// @Produce	json
// @Param	body	body		tasks.CaseInfoPayload	true	"Case information"
// @Success	200	{object}	tasks.KesslerTaskInfo
// @Failure	400	{string}	string	"Error decoding request body"
// @Failure	500	{string}	string	"Error adding task"
// @Router	/add-task/ingest/case [post]
func HandleCaseIngestAddTask(w http.ResponseWriter, r *http.Request) {
	var caseInfo tasks.CaseInfoPayload
	if err := json.NewDecoder(r.Body).Decode(&caseInfo); err != nil {
		log.Info("User Gave Bad Request for CaseIngest", "err", err)
		errorString := fmt.Sprintf("Error decoding request body: %v", err)
		http.Error(w, errorString, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	kesslerInfo, err := tasks.AddCaseTaskCastable(ctx, caseInfo)
	if err != nil {
		log.Error("Encountered Error Adding Case Task", "err", err)
		http.Error(w, fmt.Sprintf("Error adding task: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kesslerInfo)
}

// @Summary	Get Task Information
// @Description	Retrieves information about a specific task by ID
// @Tags		tasks
// @Produce	json
// @Param	id	path	string	true	"Task ID"
// @Success	200	{object}	tasks.KesslerTaskInfo
// @Failure	404	{string}	string	"Error retrieving task info"
// @Failure	501	{string}	string	"Not implemented"
// @Router	/task/{id} [get]
func HandleGetTaskInfo(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// Generic add-task helper
type castableFiling interface {
	tasks.CastableIntoFilingInfo
}

func HandleIngestAddTaskGeneric[T tasks.CastableIntoFilingInfo](w http.ResponseWriter, r *http.Request) {
	var scraperInfo T
	if err := json.NewDecoder(r.Body).Decode(&scraperInfo); err != nil {
		log.Info("User Gave Bad Request", "err", err)
		errorString := fmt.Sprintf("Error decoding request body: %v", err)
		http.Error(w, errorString, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	kesslerInfo, err := tasks.AddScraperTaskCastable(ctx, scraperInfo)
	if err != nil {
		log.Error("Encountered Error Adding Task", "err", err)
		http.Error(w, fmt.Sprintf("Error adding task: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kesslerInfo)
}

