package parexec

import (
	"fmt"
	"time"
)

func ExampleSimpleRun() {
	total := 4
	type Job struct {
		Name     string
		Duration int
	}
	jobs := make(chan *Job)
	// use buffered channel to be able to fill all jobs and start getting results without blocking workers
	results := make(chan string, total)
	work := func(stop <-chan bool) {
		for {
			select {
			case j := <-jobs:
				time.Sleep(time.Duration(j.Duration))
				results <- j.Name
			case <-stop:
				fmt.Println("closing")
				return
			}
		}
	}
	stop := SimpleRun(2, work)
	// fill jobs
	for i := 0; i < total; i++ {
		jobs <- &Job{Name: fmt.Sprint(i), Duration: i * 1000}
	}
	counter := 0
	for r := range results {
		fmt.Println(r)
		if counter++; counter >= total {
			close(stop)
			break
		}
	}
	// Wait for closing workers.
	// Normally you should do it with separate channel (in a clojure, handled in `case <-stop:`) to get acknowledges
	time.Sleep(200000)
	// Output:
	// 0
	// 1
	// 2
	// 3
	// closing
	// closing
}
