package jobs

type JobStatus string
type JobType string

const (
	Running JobStatus = "running"
	Started JobStatus = "started"
	Error   JobStatus = "error"
	Pending JobStatus = "pending"
)

// Index job types
const (
	IndexcollectionJob JobType = "index_collection"
	ReindexJob         JobType = "reindex"
	DeleteIndexJob     JobType = "delete_index"
)

// Document processing job types
const (
	Reprocess JobType = "reprocess"
)

type JobInterface interface {
	AddLog(log string)
	SetStatus(status JobStatus)
	GetStatus() JobStatus
	Start(fn interface{}, args ...interface{}) error
	SaveState() error
	LoadState() error
}

type Job struct {
	Id     string      `json:"id"`
	Status JobStatus   `json:"status"`
	Type   JobType     `json:"type"`
	Data   interface{} `json:"data"`
	JobLog []string    `json:"log"`
}
