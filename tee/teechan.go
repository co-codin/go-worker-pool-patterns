package main

import (
	"sync"
)

const (
	Fast = iota
	Slow
)

type TeeChan struct {
	chans    []chan int
	numChans int
	wgs      []*sync.WaitGroup
	wg       WaitG
}

func New(numChans int, ttype int) *TeeChan {
	chans := make([]chan int, numChans)
	wgs := make([]*sync.WaitGroup, numChans)

	for i := range numChans {
		chans[i] = make(chan int)
		wgs[i] = &sync.WaitGroup{}
	}

	var wg WaitG
	if ttype == Fast {
		wg = &WaitGStub{}
	}

	if ttype == Slow {
		wg = &WaitGNormal{}
	}

	return &TeeChan{
		numChans: numChans,
		wgs:      wgs,
		chans:    chans,
		wg:       wg,
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
				t.wg.Add(1)
				go func() {
					defer t.wgs[i].Done()
					defer t.wg.Done()
					t.chans[i] <- val
				}()

			}
			t.wg.Wait()
		}
	}()

	return t.chans
}
