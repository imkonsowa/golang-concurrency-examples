package exmaples

import (
	"fmt"
	"time"
)

// ConcurrentDigitMultiplier implementing a very simple concurrent worker patter
func ConcurrentDigitMultiplier() {
	fmt.Println(time.Now().Unix())

	n := 5
	jobs := make(chan int, n)
	results := make(chan int, n)

	// spin up 3 workers
	for w := 0; w <= 2; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= n; j++ {
		jobs <- j
	}

	close(jobs)

	for a := 1; a <= n; a++ {
		<-results
	}

	fmt.Println(time.Now().Unix())
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started job", j)
		time.Sleep(1 * time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}
