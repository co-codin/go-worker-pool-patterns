package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func worker(ctx context.Context, id int, jobs <-chan int) error {
	for j := range jobs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			fmt.Printf("Worker %d processing %d\n", id, j)
			time.Sleep(time.Second)
			if j %3 == 0 {
				return fmt.Errorf("error on job %d", j)
			}
		}
	}
	return nil
}


func main() {
	const numJobs = 5
	jobs := make(chan int, numJobs)
	var g errgroup.Group

	for w := 1; w <= 3; w++ {
		w := w
		g.Go(func() error {
			return worker(context.Background(), w, jobs)
		})
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	if err := g.Wait(); err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("All jobs completed")
    }
}