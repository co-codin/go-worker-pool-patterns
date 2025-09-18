package main

import (
	"fmt"
	"sync"
)

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

	teech := New(3, Slow)

	chans := teech.Execute(generate())

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
			fmt.Println("metrics...", v)
		}
	}()

	wg.Wait()
}
