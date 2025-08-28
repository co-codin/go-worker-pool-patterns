package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID int
}

type WorkerPool struct {
	numWorkers int
	jobQueue   chan Job
	results    chan int
	wg         sync.WaitGroup
}

func NewWorkerPool(numWorkers int, queueSize int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		jobQueue:   make(chan Job, queueSize),
		results:    make(chan int, queueSize),
	}
}
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	for job := range wp.jobQueue {
		fmt.Printf("Worker %d started job %d\n", id, job.ID)
		time.Sleep(time.Second)
		wp.results <- job.ID * 2
	}
}

func (wp *WorkerPool) Start() {
	for i := 1; i <= wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPool) AddJob(job Job) {
	wp.jobQueue <- job
}

func (wp *WorkerPool) WaitAndCollect() {
	wp.wg.Wait()
	close(wp.results)
	for result := range wp.results {
		fmt.Printf("Result: %d\n", result)
	}
}

func main() {
	numWorkers, numJobs := 3, 5
	wp := NewWorkerPool(numWorkers, numJobs)
	for i := 1; i <= numJobs; i++ {
		wp.AddJob(Job{i})
	}
	close(wp.jobQueue)
	wp.Start()
	wp.WaitAndCollect()
}
