package jobs

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

type IndexJob struct {
	Id     string      `json:"id"`
	Status JobStatus   `json:"status"`
	Type   JobType     `json:"type"`
	Data   interface{} `json:"data"`
	JobLog []string    `json:"log"`
}

func NewIndexJob(id string, jobType JobType, data interface{}) *IndexJob {
	return &IndexJob{
		Id:     id,
		Status: Pending,
		Type:   jobType,
		Data:   data,
		JobLog: []string{},
	}
}

func (j *IndexJob) AddLog(log string) {
	j.JobLog = append(j.JobLog, log)
}

func (j *IndexJob) SetStatus(status JobStatus) {
	j.Status = status
}

func (j *IndexJob) GetStatus() JobStatus {
	return j.Status
}
func (j *IndexJob) GetType() JobType {
	return j.Type
}

func (j *IndexJob) saveStateToFile() error {
	state, err := json.Marshal(j)
	if err != nil {
		log.Errorf("Error marshalling job state: %v", err)
	}
	fileName := fmt.Sprintf("job_%d_state.json", j.Id)
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf(`%s`, state))
	if err != nil {
		return err
	}

	fmt.Printf("Worker %d: state saved to %s\n", j.Id, fileName)
	return nil
}
