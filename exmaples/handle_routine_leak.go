package exmaples

import (
	"fmt"
	"time"
)

// HandleRoutineLeakExample is a borrowed example from Concurrency in go book
func HandleRoutineLeakExample() {
	doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})

		go func() {
			defer fmt.Println("doWork Exited.")
			defer close(completed)
			for {
				select {
				case s := <-strings:
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()

		return completed
	}
	done := make(chan interface{})

	completed := doWork(done, nil)

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	// block the main go routine until completed is signalled
	<-completed

	fmt.Println("Done.")
}
