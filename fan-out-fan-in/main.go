package main

import (
	"context"
	"fmt"
	"time"
)

func square(a int) int {

	return a * a
}

// CPU
func timeConsuming1() {
	counter := 0
	for range 10000000 {
		counter++
	}
}

// IO
func timeConsuming2() {
	time.Sleep(100 * time.Millisecond)
}

func timeConsuming() {
	timeConsuming2()
}

var numWorkers = 10

func generate() chan int {
	in := make(chan int)

	go func() {
		for i := range 100 {
			in <- i
		}
		close(in)
	}()

	return in
}

func main() {
	in := generate()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	now := time.Now()
	for v := range fanin(ctx, fanout(in, numWorkers, square)) {
		fmt.Println(v)
	}
	timeFanin := time.Since(now)

	fmt.Println(timeFanin)
}
