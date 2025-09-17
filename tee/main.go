package main

import (
	"fmt"
	"sync"
	"time"
)

func tee(in <- chan int,) (_, _ <- chan int) {
	out1 := make(chan int)
	out2 := make(chan int)

	go func() {
		defer close(out1)
		defer close(out2)

		for val := range in {
			wg := &sync.WaitGroup{}

			wg.Add(1)
			go func() {
				defer wg.Done()
				out1 <- val
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				out2 <- val
			}()
			
			wg.Wait()
		}
	}()

	return out1, out2
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
	ch1, ch2 := tee(generate())

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for v := range ch1 {
			fmt.Println("logging...", v)
		}
	}()

	go func() {
		defer wg.Done()
		for v := range ch2 {
			time.Sleep(time.Second)
			fmt.Println("writing...", v)
		}
	}()
	wg.Wait()
}