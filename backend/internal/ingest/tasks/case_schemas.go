package tasks

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/google/uuid"
    "github.com/hibiken/asynq"
    "kessler/internal/objects/conversations"
    "kessler/pkg/timestamp"
)

const (
    TypeIngestCase = "task:ingest_case"
)

// CaseInfoPayload represents a case and its associated filings.
// Mirrors the GenericCase Pydantic model.
type CaseInfoPayload struct {
    CaseNumber     string                 `json:"case_number"`
    CaseName       string                 `json:"case_name,omitempty"`
    CaseURL        string                 `json:"case_url,omitempty"`
    CaseType       *string                `json:"case_type,omitempty"`
    Description    *string                `json:"description,omitempty"`
    Industry       *string                `json:"industry,omitempty"`
    Petitioner     *string                `json:"petitioner,omitempty"`
    HearingOfficer *string                `json:"hearing_officer,omitempty"`
    OpenedDate     *timestamp.RFC3339Time `json:"opened_date,omitempty"`
    ClosedDate     *timestamp.RFC3339Time `json:"closed_date,omitempty"`
    Filings        []FilingInfoPayload    `json:"filings,omitempty"`
    ExtraMetadata  map[string]interface{} `json:"extra_metadata,omitempty"`
    IndexedAt      *timestamp.RFC3339Time  `json:"indexed_at,omitempty"`
}

// FilingInfoPayload represents a filing in a case.
// Mirrors the GenericFiling Pydantic model.
type FilingInfoPayload struct {
    Name          string                 `json:"name"`
    FiledDate     timestamp.RFC3339Time  `json:"filed_date"`
    PartyName     string                 `json:"party_name"`
    FilingType    string                 `json:"filing_type"`
    Description   string                 `json:"description"`
    Attachments   []AttachmentChildPayload `json:"attachments,omitempty"`
    ExtraMetadata map[string]interface{} `json:"extra_metadata,omitempty"`
}

// CastableIntoCaseInfo is implemented by types that can be converted to CaseInfoPayload.
type CastableIntoCaseInfo interface {
    IntoCaseInfo() (CaseInfoPayload, error)
}

// IntoCaseInfo returns the payload itself.
func (c CaseInfoPayload) IntoCaseInfo() (CaseInfoPayload, error) {
    return c, nil
}

// AddCaseTaskCastable enqueues a case ingestion task.
func AddCaseTaskCastable(ctx context.Context, castable CastableIntoCaseInfo) (KesslerTaskInfo, error) {
    caseInfo, err := castable.IntoCaseInfo()
    if err != nil {
        return KesslerTaskInfo{}, fmt.Errorf("error casting to CaseInfoPayload: %w", err)
    }
    task, err := NewAddCaseTask(caseInfo)
    if err != nil {
        return KesslerTaskInfo{}, fmt.Errorf("error creating case ingest task: %w", err)
    }
    client := GetClient(ctx)
    info, err := client.Enqueue(task)
    if err != nil {
        return KesslerTaskInfo{}, fmt.Errorf("error enqueueing case task: %w", err)
    }
    return KesslerTaskInfo{TaskID: info.ID, Queue: info.Queue}, nil
}

// NewAddCaseTask creates an asynq task for ingesting a case.
func NewAddCaseTask(payload CaseInfoPayload) (*asynq.Task, error) {
    data, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal case payload: %w", err)
    }
    return asynq.NewTask(TypeIngestCase, data), nil
}

// CastCaseInfoToConversation maps a case payload to ConversationInformation.
func CastCaseInfoToConversation(info CaseInfoPayload) (conversations.ConversationInformation, error) {
    if info.CaseNumber == "" {
        return conversations.ConversationInformation{}, fmt.Errorf("case_number is required")
    }
    return conversations.ConversationInformation{DocketGovID: info.CaseNumber}, nil
}
