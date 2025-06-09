package routes

import (
	"encoding/json"
	"fmt"
	"kessler/internal/ingest/tasks"
	"kessler/pkg/logger"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var log = logger.Named("ingest handler")

func DefineGlobalRouter(router *mux.Router) {
	// Version endpoint
	router.HandleFunc("/version_hash", HandleVersionHash).Methods("GET")

	// Task endpoints
	router.HandleFunc("/add-task/ingest", HandleDefaultIngestAddTask).Methods("POST")
	router.HandleFunc("/add-task/ingest/nypuc", HandleNYPUCIngestAddTask).Methods("POST")
	router.HandleFunc("/add-task/ingest/openscrapers-case", HandleCaseIngestAddTask).Methods("POST")

	router.HandleFunc("/add-task/ingest/openscrapers-caselist", HandleCaseListIngestAddTask).Methods("POST")

	// Task status endpoint
	router.HandleFunc("/task/{id}", HandleGetTaskInfo).Methods("GET")
}

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

// @Summary	Add Openscrapers Case Ingest Task
// @Description	Creates a new Openscrapers Case ingestion task
// @Tags		tasks
// @Accept	json
// @Produce	json
// @Param	body	body		tasks.OpenscrapersCaseInfoPayload	true	"Case information"
// @Success	200	{object}	tasks.KesslerTaskInfo
// @Failure	400	{string}	string	"Error decoding request body"
// @Failure	500	{string}	string	"Error adding task"
// @Router	/add-task/ingest/openscrapsers-case [post]
func HandleCaseIngestAddTask(w http.ResponseWriter, r *http.Request) {
	var caseInfo tasks.OpenscrapersCaseInfoPayload
	if err := json.NewDecoder(r.Body).Decode(&caseInfo); err != nil {
		log.Info("User Gave Bad Request for CaseIngest", zap.Error(err))
		errorString := fmt.Sprintf("Error decoding request body: %v", err)
		http.Error(w, errorString, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	kesslerInfo, err := tasks.AddCaseTaskCastable(ctx, caseInfo)
	if err != nil {
		log.Error("Encountered Error Adding Case Task", zap.Error(err))
		http.Error(w, fmt.Sprintf("Error adding task: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kesslerInfo)
}

// @Summary	Add Openscrapers CaseList Ingest Task
// @Description	Creates a new Openscrapers CaseList ingestion task
// @Tags		tasks
// @Accept	json
// @Produce	json
// @Param	body	body		[]tasks.OpenscrapersCaseListEntry	true	"Case information"
// @Success	200	{object}	string
// @Failure	400	{string}	string	"Error decoding request body"
// @Failure	500	{string}	string	"Error adding task"
// @Router	/add-task/ingest/openscrapers-caselist [post]
func HandleCaseListIngestAddTask(w http.ResponseWriter, r *http.Request) {
	var caseListInfo []tasks.OpenscrapersCaseListEntry
	if err := json.NewDecoder(r.Body).Decode(&caseListInfo); err != nil {
		log.Info("User Gave Bad Request for CaseList", zap.Error(err))
		errorString := fmt.Sprintf("Error decoding request body: %v", err)
		http.Error(w, errorString, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	for _, caseListEntry := range caseListInfo {
		caseInfo, err := caseListEntry.FetchInfoCaseInfo()
		if err != nil {
			log.Error("Error fetching info from openscrapers", zap.Error(err), zap.String("case_id", caseListEntry.CaseID))
		}
		_, err = tasks.AddCaseTaskCastable(ctx, caseInfo)
		if err != nil {
			log.Error("Encountered Error Adding Case Task", zap.Error(err), zap.String("case_id", caseListEntry.CaseID))
			// http.Error(w, fmt.Sprintf("Error adding task: %v", err), http.StatusInternalServerError)
			// return
		}
		log.Info("Successfuly finished ingest task for case", zap.String("case_id", caseListEntry.CaseID))

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("It was successful !!!!")
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
		log.Info("User Gave Bad Request", zap.Error(err))
		errorString := fmt.Sprintf("Error decoding request body: %v", err)
		http.Error(w, errorString, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	kesslerInfo, err := tasks.AddScraperFilingTaskCastable(ctx, scraperInfo)
	if err != nil {
		log.Error("Encountered Error Adding Task", zap.Error(err))
		http.Error(w, fmt.Sprintf("Error adding task: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kesslerInfo)
}
