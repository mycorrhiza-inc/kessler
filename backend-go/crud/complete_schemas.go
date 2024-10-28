package crud

import (
	"fmt"

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

type DocProcStage struct {
	PGStage         PGStage       `json:"pg_stage"`
	DocProcStatus   DocProcStatus `json:"docproc_stage"`
	IsErrored       bool          `json:"is_errored"`
	IsCompleted     bool          `json:"is_completed"`
	ErrorMsg        string        `json:"error_msg"`
	ErrorStacktrace string        `json:"error_stacktrace"`
}
type FileGeneratedExtras struct {
	Summary      string `json:"summary"`
	ShortSummary string `json:"short_summary"`
	Purpose      string `json:"purpose"`
}

type AuthorInformation struct {
	AuthorName string    `json:"author_name"`
	AuthorID   uuid.UUID `json:"author_id"`
}

type CompleteFileSchema struct {
	ID        uuid.UUID           `json:"id"`
	Extension string              `json:"extension"`
	Lang      string              `json:"lang"`
	Name      string              `json:"name"`
	Hash      string              `json:"hash"`
	IsPrivate bool                `json:"is_private"`
	Mdata     FileMetadataSchema  `json:"mdata"`
	DocTexts  []FileTextSchema    `json:"doc_texts"`
	Stage     DocProcStage        `json:"stage"`
	Extra     FileGeneratedExtras `json:"extra"`
	Authors   []AuthorInformation `json:"authors"`
}

type FileMetadataSchema struct {
	JsonObj []byte `json:"json_obj"`
}

func CompleteFileSchemaPrune(input CompleteFileSchema) FileSchema {
	return FileSchema{
		ID:        input.ID,
		Extension: input.Extension,
		Lang:      input.Lang,
		Name:      input.Name,
		Hash:      input.Hash,
		IsPrivate: input.IsPrivate,
	}
}

func FileSchemaToComplete(input FileSchema) (CompleteFileSchema, error) {
	return_schema := CompleteFileSchema{
		ID:        input.ID,
		Extension: input.Extension,
		Lang:      input.Lang,
		Name:      input.Name,
		Hash:      input.Hash,
		IsPrivate: input.IsPrivate,
	}
	// TODO: Query Metadata json and also get other stuff
	return return_schema, fmt.Errorf("not implemented")
}

func ConvertToCreationData(updateInfo CompleteFileSchema) FileCreationDataRaw {
	creationData := FileCreationDataRaw{
		Extension: pgtype.Text{String: updateInfo.Extension, Valid: true},
		Lang:      pgtype.Text{String: updateInfo.Lang, Valid: true},
		Name:      pgtype.Text{String: updateInfo.Name, Valid: true},
		Hash:      pgtype.Text{String: updateInfo.Hash, Valid: true},
	}
	return creationData
}
