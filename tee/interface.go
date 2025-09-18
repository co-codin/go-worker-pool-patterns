package main

type WaitG interface {
	Add(int)
	Done()
	Wait()
}
