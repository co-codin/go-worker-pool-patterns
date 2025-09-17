package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func tee(ctx context.Context, in <-chan int, numChans int) []chan int {
	chans := make([]chan int, numChans)

	for i := range numChans {
		chans[i] = make(chan int)
	}

	go func() {
		for i := range numChans {
			defer close(chans[i])
		}

		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-in:
				if !ok {
					return
				}

				wg := &sync.WaitGroup{}

				for i := range numChans {
					wg.Add(1)
					go func() {
						defer wg.Done()

						select {
						case chans[i] <- val:
						case <-ctx.Done():
							return
						}

					}()
				}
				wg.Wait()
			}
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	chans := tee(ctx, generate(), 2)

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
