package jobs

import "sync"

type JobManager struct {
	mu           sync.Mutex             // Protects access to jobs and state
	running_jobs map[string]*Job        // Map of job IDs to Job structs
	jobQueue     chan *Job              // Channel for managing a pool of jobs
	state        map[string]interface{} // Persistent state of jobs
	done         chan struct{}          // Signal to stop all jobs
}

func NewJobManager() *JobManager {
	return &JobManager{
		running_jobs: make(map[string]*Job),
		jobQueue:     make(chan *Job),
		state:        make(map[string]interface{}),
		done:         make(chan struct{}),
	}
}

// func (m *JobManager) Start(fn interface{}, args ...interface{}) {
// 	go func() {
// 		for {
// 			select {
// 			case job := <-m.jobQueue:
// 				m.mu.Lock()
// 				m.running_jobs[job.Id] = job
// 				m.mu.Unlock()
// 				go job.Run()
// 			case <-m.done:
// 				return
// 			}
// 		}
// 	}()
// }
