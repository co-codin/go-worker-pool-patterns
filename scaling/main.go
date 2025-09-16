package main

import "sync"


type DynamicPool struct {
	minWorkers int
	maxWorkers int
	jobQueue chan int
	wg sync.WaitGroup
	mu sync.Mutex
	active int
}

func NewDynamicPool(min, max int) *DynamicPool {
	dp := &DynamicPool{minWorkers: min, maxWorkers: max, jobQueue: make(chan int, 100)}
	for i := 0; i < min; i++ {
		dp.startWorker()
	}
	return dp
}

func (dp *DynamicPool) startWorker() {
	
}