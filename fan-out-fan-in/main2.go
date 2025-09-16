package main

// import (
// 	"fmt"
// 	"time"
// 	"sync"
// )

// //Fans tasks out to multiple workers for parallel processing,
// //then fans results in to a single collector.
// //Great for independent, parallelizable tasks like batch API requests.

// type Task struct {
// 	ID int
// }

// func processTask(task Task) string {
// 	time.Sleep(time.Second)
// 	return fmt.Sprintf("Processed task %d", task.ID)
// }

// func worker(tasks <- chan Task, results chan<- string, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for task := range tasks {
// 		results <- processTask(task)
// 	}
// }

// func main() {
// 	tasks := make(chan Task, 5)
// 	results := make(chan string, 5)

// 	var wg sync.WaitGroup

// 	wg.Add(3)

// 	for i := 0; i < 3; i++ {
// 		go worker(tasks, results, &wg)
// 	}

// 	go func() {
// 		for i := 1; i <= 5; i++ {
// 			tasks <- Task{ID: i}
// 		}
// 		close(tasks)
// 	}()

// 	go func() {
// 		wg.Wait()
// 		close(results)
// 	}()

// 	for result := range results {
// 		fmt.Println(result)
// 	}
// }
