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

type FileChildTextSource struct {
	IsOriginalText bool   `json:"is_original_text"`
	Text           string `json:"text"`
	Language       string `json:"language"`
}

func ChildTextSouceToRealTextSource(child_source FileChildTextSource, id uuid.UUID) FileTextSchema {
	return FileTextSchema{
		FileID:         pgtype.UUID{Bytes: id, Valid: true},
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

type AuthorInformation struct {
	AuthorName      string    `json:"author_name"`
	IsPerson        bool      `json:"is_person"`
	IsPrimaryAuthor bool      `json:"is_primary_author"`
	AuthorID        uuid.UUID `json:"author_id"`
}

type JuristictionInformation struct {
	Country        string                 `json:"country"`
	State          string                 `json:"state"`
	Municipality   string                 `json:"municipality"`
	Agency         string                 `json:"agency"`
	ProceedingName string                 `json:"proceeding_name"`
	ExtraObject    map[string]interface{} `json:"extra_object"`
}

type CompleteFileSchema struct {
	ID           uuid.UUID               `json:"id"`
	Verified     bool                    `json:"verified"`
	Extension    string                  `json:"extension"`
	Lang         string                  `json:"lang"`
	Name         string                  `json:"name"`
	Hash         string                  `json:"hash"`
	IsPrivate    bool                    `json:"is_private"`
	Mdata        FileMetadataSchema      `json:"mdata"`
	Stage        DocProcStage            `json:"stage"`
	Extra        FileGeneratedExtras     `json:"extra"`
	Authors      []AuthorInformation     `json:"authors"`
	Juristiction JuristictionInformation `json:"juristiction"`
	// To heavy to include in a default file schema unless the user specifies they want it
	DocTexts []FileChildTextSource `json:"doc_texts"`
}

type SemiCompleteFileSchema struct {
	ID           uuid.UUID               `json:"id"`
	Verified     bool                    `json:"verified"`
	Extension    string                  `json:"extension"`
	Lang         string                  `json:"lang"`
	Name         string                  `json:"name"`
	Hash         string                  `json:"hash"`
	IsPrivate    bool                    `json:"is_private"`
	Mdata        FileMetadataSchema      `json:"mdata"`
	Stage        DocProcStage            `json:"stage"`
	Extra        FileGeneratedExtras     `json:"extra"`
	Authors      []AuthorInformation     `json:"authors"`
	Juristiction JuristictionInformation `json:"juristiction"`
}
type FileMetadataSchema map[string]interface{}

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
		Verified:  pgtype.Bool{Bool: updateInfo.Verified, Valid: true},
	}
	return creationData
}
