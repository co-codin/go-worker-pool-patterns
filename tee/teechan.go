package main

import (
	"sync"
)

type TeeChan struct {
	chans    []chan int
	numChans int
	wgs      []*sync.WaitGroup
}

func New(numChans int) *TeeChan {
	chans := make([]chan int, numChans)
	wgs := make([]*sync.WaitGroup, numChans)

	for i := range numChans {
		chans[i] = make(chan int)
		wgs[i] = &sync.WaitGroup{}
	}

	return &TeeChan{
		numChans: numChans,
		wgs:      wgs,
		chans:    chans,
	}
}

func (t *TeeChan) Execute(in chan int) []chan int {
	go func() {
		defer func() {
			for i := range t.numChans {
				go func() {
					t.wgs[i].Wait()
					close(t.chans[i])
				}()
			}

		}()

		for val := range in {
			for i := range t.numChans {
				t.wgs[i].Add(1)
				go func() {
					defer t.wgs[i].Done()
					t.chans[i] <- val
				}()

			}
		}
	}()

	return t.chans
}
