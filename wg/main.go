package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func processData(ctx context.Context, v int) int {
	ch := make(chan struct{})

	go func() {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		close(ch)
	}()

	select {
	case <- ch:
		return v *2
	case <- ctx.Done():
		return 0
	}
}

func main() {
	in := make(chan int)
	out := make(chan int)

	go func() {
		for i := range 10 {
			in <- i
		}
		close(in)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	start := time.Now()
	processParallel(ctx, in, out, 5)

	for v := range out {
		fmt.Println("v = ", v)
	}
	fmt.Println("duration:", time.Since(start))

}

func processParallel(ctx context.Context, in <-chan int, out chan<- int, numWorkers int) {
	wg := &sync.WaitGroup{}
	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range in {
				select {
				case out <- processData(ctx, v):
				case <- ctx.Done():
					fmt.Println("timeout")
				}
				
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()
}
