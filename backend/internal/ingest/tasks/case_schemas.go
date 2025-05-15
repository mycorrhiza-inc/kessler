package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/internal/objects/conversations"
	"kessler/pkg/hashes"
	"kessler/pkg/timestamp"

	"github.com/hibiken/asynq"
)

type KesslerTaskInfo struct {
	TaskID string `json:"task_id"`
	Queue  string `json:"queue"`
	State  string `json:"state"`
	Status string `json:"status"`
}

func GenerateTaskInfoFromInfo(info asynq.TaskInfo) KesslerTaskInfo {
	return KesslerTaskInfo{
		TaskID: info.ID,
		Queue:  info.Queue,
	}
}

const (
	TypeIngestCase = "task:ingest_case"
)

// OpenscrapersCaseInfoPayload represents a case and its associated filings.
// Mirrors the GenericCase Pydantic model.
type OpenscrapersCaseInfoPayload struct {
	CaseNumber     string                 `json:"case_number"`
	CaseName       string                 `json:"case_name,omitempty"`
	CaseURL        string                 `json:"case_url,omitempty"`
	CaseType       string                 `json:"case_type,omitempty"`
	Description    string                 `json:"description,omitempty"`
	Industry       string                 `json:"industry,omitempty"`
	Petitioner     string                 `json:"petitioner,omitempty"`
	HearingOfficer string                 `json:"hearing_officer,omitempty"`
	OpenedDate     timestamp.RFC3339Time  `json:"opened_date,omitempty"`
	ExtraMetadata  map[string]interface{} `json:"extra_metadata,omitempty"`
	IndexedAt      timestamp.RFC3339Time  `json:"indexed_at,omitempty"`
	Filings        []FilingChildInfo      `json:"filings,omitempty"`
}

type CaseInfoMinimal struct {
	CaseNumber     string                 `json:"case_number"`
	CaseName       string                 `json:"case_name,omitempty"`
	CaseURL        string                 `json:"case_url,omitempty"`
	CaseType       string                 `json:"case_type,omitempty"`
	Description    string                 `json:"description,omitempty"`
	Industry       string                 `json:"industry,omitempty"`
	Petitioner     string                 `json:"petitioner,omitempty"`
	HearingOfficer string                 `json:"hearing_officer,omitempty"`
	OpenedDate     timestamp.RFC3339Time  `json:"opened_date,omitempty"`
	ExtraMetadata  map[string]interface{} `json:"extra_metadata,omitempty"`
	IndexedAt      timestamp.RFC3339Time  `json:"indexed_at,omitempty"`
}

func (c OpenscrapersCaseInfoPayload) IntoCaseInfoMinimal() CaseInfoMinimal {
	return CaseInfoMinimal{
		CaseNumber:     c.CaseNumber,
		CaseName:       c.CaseName,
		CaseURL:        c.CaseURL,
		CaseType:       c.CaseType,
		Description:    c.Description,
		Industry:       c.Industry,
		Petitioner:     c.Petitioner,
		HearingOfficer: c.HearingOfficer,
		OpenedDate:     c.OpenedDate,
		ExtraMetadata:  c.ExtraMetadata,
		IndexedAt:      c.IndexedAt,
	}
}

// FilingChildInfo represents a filing in a case.
// Mirrors the GenericFiling Pydantic model.
type FilingChildInfo struct {
	Name          string                 `json:"name"`
	FiledDate     timestamp.RFC3339Time  `json:"filed_date"`
	PartyName     string                 `json:"party_name"`
	FilingType    string                 `json:"filing_type"`
	Description   string                 `json:"description"`
	Attachments   []AttachmentChildInfo  `json:"attachments,omitempty"`
	ExtraMetadata map[string]interface{} `json:"extra_metadata,omitempty"`
}

type AttachmentChildInfo struct {
	Lang          string             `json:"lang"`
	Name          string             `json:"name"`
	Extension     string             `json:"extension"`
	URL           string             `json:"url"`
	Hash          hashes.KesslerHash `json:"hash"`
	Mdata         map[string]any     `json:"mdata"`
	RawAttachment RawAttachmentData  `json:"raw_attachment"`
}

type RawAttachmentData struct {
	Hash             string `json:"hash"`
	Name             string `json:"name"`
	Extension        string `json:"extension"`
	GetAttachmentUrl string `json:"get_attachment_url"`
	TextObjects      []struct {
		Quality   int    `json:"quality"`
		Language  string `json:"language"`
		Text      string `json:"text"`
		Timestamp string `json:"timestamp"`
	} `json:"text_objects"`
}

// CastableIntoCaseInfo is implemented by types that can be converted to CaseInfoPayload.
type CastableIntoCaseInfo interface {
	IntoCaseInfo() (OpenscrapersCaseInfoPayload, error)
}

// IntoCaseInfo returns the payload itself.
func (c OpenscrapersCaseInfoPayload) IntoCaseInfo() (OpenscrapersCaseInfoPayload, error) {
	return c, nil
}

// AddCaseTaskCastable enqueues a case ingestion task.
func AddCaseTaskCastable(ctx context.Context, castable CastableIntoCaseInfo) (KesslerTaskInfo, error) {
	caseInfo, err := castable.IntoCaseInfo()
	if err != nil {
		return KesslerTaskInfo{}, fmt.Errorf("error casting to CaseInfoPayload: %w", err)
	}
	err = IngestOpenscrapersCase(ctx, caseInfo)
	if err != nil {
		return KesslerTaskInfo{}, err
	}

	return KesslerTaskInfo{}, nil
	// task, err := NewAddCaseTask(caseInfo)
	// if err != nil {
	// 	return KesslerTaskInfo{}, fmt.Errorf("error creating case ingest task: %w", err)
	// }
	// client := GetClient(ctx)
	// info, err := client.Enqueue(task)
	// if err != nil {
	// 	return KesslerTaskInfo{}, fmt.Errorf("error enqueueing case task: %w", err)
	// }
	// return KesslerTaskInfo{TaskID: info.ID, Queue: info.Queue}, nil
}

// NewAddCaseTask creates an asynq task for ingesting a case.
func NewAddCaseTask(payload OpenscrapersCaseInfoPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal case payload: %w", err)
	}
	return asynq.NewTask(TypeIngestCase, data), nil
}

// CastCaseInfoToConversation maps a case payload to ConversationInformation.
func CastCaseInfoToConversation(info OpenscrapersCaseInfoPayload) (conversations.ConversationInformation, error) {
	if info.CaseNumber == "" {
		return conversations.ConversationInformation{}, fmt.Errorf("case_number is required")
	}
	return conversations.ConversationInformation{DocketGovID: info.CaseNumber}, nil
}
