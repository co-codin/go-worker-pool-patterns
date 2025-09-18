package main



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


	teech := New(3)

	chans := teech.Execute(generate())
}
