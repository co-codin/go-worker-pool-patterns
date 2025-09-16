package main

import (
	"fmt"
	"sync"
	"time"
)

type DynamicPool struct {
	minWorkers int
	maxWorkers int
	jobQueue   chan int
	wg         sync.WaitGroup
	mu         sync.Mutex
	active     int
}

func NewDynamicPool(min, max int) *DynamicPool {
	dp := &DynamicPool{minWorkers: min, maxWorkers: max, jobQueue: make(chan int, 100)}
	for i := 0; i < min; i++ {
		dp.startWorker()
	}
	return dp
}

func (dp *DynamicPool) startWorker() {
	dp.mu.Lock()
	dp.active++
	dp.mu.Unlock()
	dp.wg.Add(1)
	go func() {
		defer dp.wg.Done()
		defer func() {
			dp.mu.Lock()
			dp.active--
			dp.mu.Unlock()
		}()
		for job := range dp.jobQueue {
			fmt.Printf("Processing job %d\n", job)
			time.Sleep(time.Second)
		}
	}()
}

func (dp *DynamicPool) AddJob(job int) {
	dp.jobQueue <- job
	dp.mu.Lock()
	if len(dp.jobQueue) > dp.active && dp.active < dp.maxWorkers {
		dp.startWorker()
	}
	dp.mu.Unlock()
}


func (dp *DynamicPool) Shutdown() {
	close(dp.jobQueue)
	dp.wg.Wait()
}

func main() {
	dp := NewDynamicPool(2, 5)
	for i := 1; i <= 10; i++ {
		dp.AddJob(1)
	}
	dp.Shutdown()
}