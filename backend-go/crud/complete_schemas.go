package crud

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type PGStage string

const (
	StageStatePending    PGStage = "pending"
	StageStateProcessing PGStage = "processing"
	StageStateCompleted  PGStage = "completed"
	StageStateErrored    PGStage = "errored"
)

type DocProcStatus string

const (
	Unprocessed            DocProcStatus = "unprocessed"
	Completed              DocProcStatus = "completed"
	EncountersAnalyzed     DocProcStatus = "encounters_analyzed"
	OrganizationAssigned   DocProcStatus = "organization_assigned"
	SummarizationCompleted DocProcStatus = "summarization_completed"
	EmbeddingsCompleted    DocProcStatus = "embeddings_completed"
	UploadDocumentToDB     DocProcStatus = "upload_document_to_db"
	Stage3                 DocProcStatus = "stage3"
	Stage2                 DocProcStatus = "stage2"
	Stage1                 DocProcStatus = "stage1"
)

type FileChildTextSource struct {
	IsOriginalText bool   `json:"is_original_text"`
	Text           string `json:"text"`
	Language       string `json:"language"`
}

func ChildTextSouceToRealTextSource(child_source FileChildTextSource, id uuid.UUID) FileTextSchema {
	return FileTextSchema{
		FileID:         id,
		IsOriginalText: child_source.IsOriginalText,
		Text:           child_source.Text,
		Language:       child_source.Language,
	}
}

type DocProcStage struct {
	PGStage            PGStage       `json:"pg_stage"`
	DocProcStatus      DocProcStatus `json:"docproc_stage"`
	IsErrored          bool          `json:"is_errored"`
	IsCompleted        bool          `json:"is_completed"`
	ProcessingErrorMsg string        `json:"processing_error_msg"`
	DatabaseErrorMsg   string        `json:"database_error_msg"`
}
type FileGeneratedExtras struct {
	Summary        string  `json:"summary"`
	ShortSummary   string  `json:"short_summary"`
	Purpose        string  `json:"purpose"`
	Impressiveness float64 `json:"impressiveness"`
}

type JuristictionInformation struct {
	Country        string                 `json:"country"`
	State          string                 `json:"state"`
	Municipality   string                 `json:"municipality"`
	Agency         string                 `json:"agency"`
	ProceedingName string                 `json:"proceeding_name"`
	ExtraObject    map[string]interface{} `json:"extra_object"`
}

type ConversationInformation struct {
	ID          uuid.UUID `json:"id"`
	DocketID    string    `json:"docket_id"`
	State       string    `json:"state"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type AuthorInformation struct {
	AuthorName      string    `json:"author_name"`
	IsPerson        bool      `json:"is_person"`
	IsPrimaryAuthor bool      `json:"is_primary_author"`
	AuthorID        uuid.UUID `json:"author_id"`
}

// Organization Schemas

type OrganizationSchemaComplete struct {
	ID               uuid.UUID    `json:"id"`
	Name             string       `json:"name"`
	FilesAuthored    []FileSchema `json:"files_authored"`
	FilesAuthoredIDs []uuid.UUID  `json:"files_authored_ids"`
}

func CompleteFileSchemaPrune(input CompleteFileSchema) FileSchema {
	return FileSchema{
		ID:        input.ID,
		Verified:  input.Verified,
		Extension: input.Extension,
		Lang:      input.Lang,
		Name:      input.Name,
		Hash:      input.Hash,
		IsPrivate: input.IsPrivate,
	}
}

func CompleteFileSchemaInflateFromPartialSchema(input FileSchema) CompleteFileSchema {
	return_schema := CompleteFileSchema{
		ID:        input.ID,
		Verified:  input.Verified,
		Extension: input.Extension,
		Lang:      input.Lang,
		Name:      input.Name,
		Hash:      input.Hash,
		IsPrivate: input.IsPrivate,
	}
	// TODO: Query Metadata json and also get other stuff
	return return_schema
}

func ConvertToCreationData(updateInfo CompleteFileSchema) FileCreationDataRaw {
	creationData := FileCreationDataRaw{
		Extension: updateInfo.Extension,
		Lang:      updateInfo.Lang,
		Name:      updateInfo.Name,
		Hash:      updateInfo.Hash,
		Verified:  pgtype.Bool{Bool: updateInfo.Verified, Valid: true},
	}
	return creationData
}
