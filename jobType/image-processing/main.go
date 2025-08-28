package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

type Job struct {
	ID       int
	ImageUrl string
	Size     int
}

type Result struct {
	JobID     int
	NewSize   int
	Error     error
	TimeSpent time.Duration
}

func imageProcessor(jobs <-chan Job, results chan<- Result) {
	for job := range jobs {
		start := time.Now()
		err := ProcessImage(job)
		results <- Result{
			JobID:     job.ID,
			NewSize:   job.Size / 2,
			Error:     err,
			TimeSpent: time.Since(start),
		}
	}
}

func ProcessImage(job Job) error {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	if rand.Float32() < 0.1 {
		return fmt.Errorf("image processing failed for %s", job.ImageUrl)
	}
	return nil
}

func main() {
	numCPU := runtime.NumCPU()
	numWorkers := numCPU * 2
	const numJobs = 5

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	for w := 1; w <= numWorkers; w++ {
		go imageProcessor(jobs, results)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- Job{
			ID:       j,
			ImageUrl: fmt.Sprintf("image%d.jpg", j),
			Size:     100 * j,
		}
	}

	close(jobs)

	for a := 1; a <= numJobs; a++ {
		res := <-results
		if res.Error != nil {
			fmt.Printf("Error on job %d: %v\n", res.JobID, res.Error)
		} else {
			fmt.Printf("Processed job %d in %v\n", res.JobID, res.TimeSpent)
		}
	}
}
