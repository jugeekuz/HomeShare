package job

import (
	"sync"
	"time"
)

type Job struct {
	mu         sync.Mutex
	jobId	   string
	lastAccess time.Time
}

type JobManager struct {
	jobs      map[string]*Job
	mapMu     sync.RWMutex // To protect against race conditions in accessing the jobs map
	timeout   time.Duration
	closeChan chan struct{}
}

func NewJobManager(timeout time.Duration) *JobManager {
	jm := &JobManager{
		jobs:      make(map[string]*Job),
		timeout:   timeout,
		closeChan: make(chan struct{}),
	}
	go jm.cleanupStaleJobs()
	return jm
}

func (jm *JobManager) AcquireJob(jobId string) (bool) {
	jm.mapMu.Lock()
	defer jm.mapMu.Unlock()

	if job, exists := jm.jobs[jobId]; exists {
		job.lastAccess = time.Now()
		return false // Job is being processed
	}

	job := &Job{
		jobId:       	jobId,
		lastAccess: 	time.Now(),
	}
	job.mu.Lock()
	jm.jobs[jobId] = job
	return true
}

func (jm *JobManager) ReleaseJob(jobId string) {
	jm.mapMu.Lock()
	defer jm.mapMu.Unlock()

	if job, exists := jm.jobs[jobId]; exists {
		job.mu.Unlock()
		delete(jm.jobs, jobId)
	}
}

func (jm *JobManager) cleanupStaleJobs() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			jm.mapMu.Lock()
			for id, job := range jm.jobs {
				if time.Since(job.lastAccess) > jm.timeout {
					delete(jm.jobs, id)
				}
			}
			jm.mapMu.Unlock()
		case <-jm.closeChan:
			return
		}
	}
}

func (jm *JobManager) Close() {
	close(jm.closeChan)
	jm.mapMu.Lock()
	defer jm.mapMu.Unlock()
	for id, _ := range jm.jobs {
		delete(jm.jobs, id)
	}
}