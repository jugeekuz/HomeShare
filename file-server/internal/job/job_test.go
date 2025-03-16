package job

import (
	"sync"
	"testing"
	"time"
)

func TestJobManager(t *testing.T) {
	t.Run("Test acquire & release job", func(t *testing.T) {
		jm := NewJobManager(5 * time.Minute)
		jobID := "job1"

		if acquired := jm.AcquireJob(jobID); !acquired {
			t.Errorf("expected job %s to be acquired, got false", jobID)
		}

		if acquired := jm.AcquireJob(jobID); acquired {
			t.Errorf("expected job %s to not be acquired twice, got true", jobID)
		}

		jm.ReleaseJob(jobID)

		if acquired := jm.AcquireJob(jobID); !acquired {
			t.Errorf("expected job %s to be acquired after release, got false", jobID)
		}

		jm.ReleaseJob(jobID)
		jm.Close()
	})

	t.Run("Test Concurrency", func(t *testing.T) {
		jm := NewJobManager(5 * time.Minute)
		jobID := "concurrentJob"
		var wg sync.WaitGroup
		jobsCount := 0
		attempts := 10
		var mu sync.Mutex

		for i := 0; i < attempts; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if jm.AcquireJob(jobID) {
					mu.Lock()
					jobsCount++
					mu.Unlock()
				}
			}()
		}
		wg.Wait()

		if jobsCount != 1 {
			t.Errorf("expected exactly one successful acquisition, got %d", jobsCount)
		}
		jm.ReleaseJob(jobID)
		jm.Close()
	})

	t.Run("Test Cleanup Stale jobs", func(t *testing.T) {
		jm := &JobManager{
			jobs:      make(map[string]*Job),
			timeout:   50 * time.Millisecond,
			closeChan: make(chan struct{}),
		}

		jobID := "staleJob"

		if acquired := jm.AcquireJob(jobID); !acquired {
			t.Errorf("expected job %s to be acquired, got false", jobID)
		}

		go jm.cleanupStaleJobs(5 * time.Millisecond)

		time.Sleep(60 * time.Millisecond)

		jm.mapMu.Lock()
		if _, exists := jm.jobs[jobID]; exists {
			t.Errorf("expected job %s to be cleaned up as stale", jobID)
		}
		jm.mapMu.Unlock()

		jm.Close()
	})

	t.Run("Test Close Method", func(t *testing.T) {
		jm := NewJobManager(5 * time.Minute)
		jobIDs := []string{"job1", "job2", "job3"}

		for _, id := range jobIDs {
			jm.AcquireJob(id)
		}

		jm.Close()

		jm.mapMu.RLock()
		if len(jm.jobs) != 0 {
			t.Errorf("expected jobs map to be empty after Close, but found %d", len(jm.jobs))
		}
		jm.mapMu.RUnlock()
	})
}
