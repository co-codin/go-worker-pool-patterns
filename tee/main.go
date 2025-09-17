package main

import (
	"fmt"
	"sync"
)

func tee(in <-chan int, numChans int) []chan int {
	chans := make([]chan int, numChans)

	for i := range numChans {
		chans[i] = make(chan int)
	}

	go func() {
		for i := range numChans {
			defer close(chans[i])
		}

		for val := range in {
			wg := &sync.WaitGroup{}

			for i := range numChans {
				wg.Add(1)
				go func() {
					defer wg.Done()
					chans[i] <- val
				}()
			}
			wg.Wait()
		}
	}()

	return chans
}

func generate() chan int {
	ch := make(chan int)

	go func() {
		for i := range 5 {
			ch <- i
		}
		close(ch)
	}()

	return ch
}

func main() {
	chans := tee(generate(), 2)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for v := range chans[0] {
			fmt.Println("logging...", v)
		}
	}()

	go func() {
		defer wg.Done()
		for v := range chans[1] {
			fmt.Println("logging...", v)
		}
	}()

	wg.Wait()
}
