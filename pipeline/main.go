package main

import (
	"fmt"
)

func generator(input []int) chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, data := range input {
			ch <- data
		}
	}()
	return ch
}

func add(inputCh chan int) chan int {
	res := make(chan int)
	go func() {
		defer close(res)
		for data := range inputCh {
			res <- data + 1
		}
	}()

	return res
}

func multiply(inputCh chan int) chan int {
	res := make(chan int)
    go func() {
        defer close(res)
        for data := range inputCh {
            res <- data * 2
        }
    }()
    return res
}

func main() {
	input := []int{1, 2, 3}
    inputCh := generator(input)
    resultCh := multiply(add(inputCh))
    for res := range resultCh {
        fmt.Println(res)
    }
}