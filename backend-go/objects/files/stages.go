package files

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
	PGStage            PGStage       `json:"pg_stage"`
	DocProcStatus      DocProcStatus `json:"docproc_stage"`
	SkipProcessing     bool          `json:"skip_processing"`
	IsErrored          bool          `json:"is_errored"`
	IsCompleted        bool          `json:"is_completed"`
	IngestErrorMsg     string        `json:"ingest_error_msg"`
	ProcessingErrorMsg string        `json:"processing_error_msg"`
	DatabaseErrorMsg   string        `json:"database_error_msg"`
}
