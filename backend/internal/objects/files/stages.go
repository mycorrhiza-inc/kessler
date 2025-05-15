package files

type PGStage string

const (
	PGStagePending    PGStage = "pending"
	PGStageProcessing PGStage = "processing"
	PGStageCompleted  PGStage = "completed"
	PGStageErrored    PGStage = "errored"
)

type DocProcStatus string

const (
	DocStatusUnprocessed            DocProcStatus = "unprocessed"
	DocStatusCompleted              DocProcStatus = "completed"
	DocStatusEncountersAnalyzed     DocProcStatus = "encounters_analyzed"
	DocStatusOrganizationAssigned   DocProcStatus = "organization_assigned"
	DocStatusSummarizationCompleted DocProcStatus = "summarization_completed"
	DocStatusEmbeddingsCompleted    DocProcStatus = "embeddings_completed"
	DocStatusUploadDocumentToDB     DocProcStatus = "upload_document_to_db"
	DocStatusTextCompleted          DocProcStatus = "text_completed"
	DocStatusRawTextCompleted       DocProcStatus = "raw_text_completed"
	DocStatusBeginProcessing        DocProcStatus = "begin_processing"
)

func (status DocProcStatus) Index() int {
	switch status {
	case DocStatusUnprocessed:
		return 0
	case DocStatusBeginProcessing:
		return 1
	case DocStatusRawTextCompleted:
		return 2
	case DocStatusTextCompleted:
		return 3
	case DocStatusEncountersAnalyzed:
		return 4
	case DocStatusOrganizationAssigned:
		return 5
	case DocStatusSummarizationCompleted:
		return 6
	case DocStatusEmbeddingsCompleted:
		return 7
	case DocStatusUploadDocumentToDB:
		return 8
	case DocStatusCompleted:
		return 9
	default:
		return -1
	}
}

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
