package jobs

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

type IndexJob struct {
	Job
}

func NewIndexJob(id string, jobType JobType, data interface{}) *IndexJob {
	return &IndexJob{
		Job: Job{
			Id:     id,
			Status: Pending,
			JobLog: []string{},
		},
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

func (j *IndexJob) SaveStateToDisk() error {
	state, err := json.Marshal(j)
	if err != nil {
		log.Errorf("Error marshalling job state: %v", err)
	}
	fileName := fmt.Sprintf("job_%s_state.json", j.Id)
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(state)
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Worker %s: state saved to %s\n", j.Id, fileName))
	return nil
}

func (j *IndexJob) LoadStateFromDisk() error {
	fileName := fmt.Sprintf("job_%s_state.json", j.Id)
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	state := &IndexJob{}
	err = json.NewDecoder(file).Decode(state)
	if err != nil {
		return err
	}

	j.Status = state.Status
	j.JobLog = state.JobLog

	log.Info(fmt.Sprintf("Worker %s: state loaded from %s\n", j.Id, fileName))
	return nil
}

// func (j *IndexJob) Start(f func() error) {
// 	j.func ()  {

// // 	} = f
// }

func IndexOrganization(id string, job *IndexJob) error {
	// get organization data from postgres

	// index organization data in quickwit

	return nil
}
