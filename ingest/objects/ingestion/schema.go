package ingestion

import (
	"time"

	"github.com/google/uuid"
)

type TaskType string

const (
	TaskTypeAddFileScraper             TaskType = "add_file_scraper"
	TaskTypeProcessExistingFile        TaskType = "process_existing_file"
	TaskTypeAddNypucConversationDocket TaskType = "add_nypuc_conversation_docket"
)

type DatabaseInteraction string

const (
	DatabaseInteractionNone              DatabaseInteraction = "none"
	DatabaseInteractionInsertLater       DatabaseInteraction = "insert_later"
	DatabaseInteractionUpdate            DatabaseInteraction = "update"
	DatabaseInteractionInsert            DatabaseInteraction = "insert"
	DatabaseInteractionInsertReportLater DatabaseInteraction = "insert_report_later"
	DatabaseInteractionInsertReport      DatabaseInteraction = "insert_report"
	DatabaseInteractionUpdateReport      DatabaseInteraction = "update_report"
)

type ScraperInfo struct {
	FileURL               string `json:"file_url"`
	Text                  string `json:"text"`
	Hash                  string `json:"hash"`
	Name                  string `json:"name"`
	PublishedDate         string `json:"published_date"`
	InternalSourceName    string `json:"internal_source_name"`
	DocketID              string `json:"docket_id"`
	State                 string `json:"state"`
	AuthorIndividual      string `json:"author_individual"`
	AuthorIndividualEmail string `json:"author_individual_email"`
	AuthorOrganisation    string `json:"author_organisation"`
	FileClass             string `json:"file_class"`
	FileType              string `json:"file_type"`
	Lang                  string `json:"lang"`
	ItemNumber            string `json:"item_number"`
}

type BulkProcessInfo struct {
	GenerateReport      bool                `json:"generate_report"`
	ReportID            string              `json:"report_id"`
	DatabaseInteraction DatabaseInteraction `json:"database_interaction"`
	OverrideScrapeInfo  ScraperInfo         `json:"override_scrape_info"`
}

type BulkProcessSchema struct {
	ScraperInfoList []ScraperInfo   `json:"scraper_info_list"`
	BulkInfo        BulkProcessInfo `json:"bulk_info"`
}

type Task struct {
	ID               uuid.UUID              `json:"id"`
	Priority         bool                   `json:"priority"`
	DatabaseInteract DatabaseInteraction    `json:"database_interact"`
	TaskType         TaskType               `json:"task_type"`
	TableName        string                 `json:"table_name"`
	Kwargs           map[string]interface{} `json:"kwargs"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	Error            string                 `json:"error"`
	Completed        bool                   `json:"completed"`
	Success          bool                   `json:"success"`
	FollowupTaskID   *uuid.UUID             `json:"followup_task_id,omitempty"`
	FollowupTaskURL  *string                `json:"followup_task_url,omitempty"`
	Obj              any                    `json:"obj"`
}

func NewTask(
	obj any,
	databaseInteraction DatabaseInteraction,
	priority bool,
	kwargs map[string]interface{},
	taskType TaskType,
) *Task {
	task := &Task{
		ID:               uuid.New(),
		Priority:         priority,
		DatabaseInteract: databaseInteraction,
		TaskType:         taskType,
		TableName:        "",
		Kwargs:           kwargs,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Error:            "",
		Completed:        false,
		Success:          false,
		Obj:              obj,
	}
	return task
}

func (t Task) URL() string {
	return "https://api.kessler.xyz/ingest_v1/tasks/" + t.ID.String()
}
