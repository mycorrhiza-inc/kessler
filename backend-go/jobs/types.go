package jobs

type JobStatus string
type JobType string

const (
	Running JobStatus = "running"
	Started JobStatus = "started"
	Error   JobStatus = "error"
	Pending JobStatus = "pending"
)

const (
	IndexcollectionJob JobType = "index_collection"
	ReindexJob         JobType = "reindex"
	DeleteIndexJob     JobType = "delete_index"
)

type Job interface {
	AddLog(log string)
	SetStatus(status JobStatus)
	GetStatus() JobStatus
	GetType() JobType
	SaveState() error
	LoadState() error
}
