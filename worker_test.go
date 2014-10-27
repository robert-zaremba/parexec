package parexec

import (
	"fmt"
	"testing"
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

func ExampleSimpleRun_singleloop() {
	var jobs = make(chan int)
	var results = make(chan int, 2)
	work := func(stop <-chan bool) {
		var j int
		for {
			select {
			case j = <-jobs:
				results <- j
			case <-stop:
				return
			}
		}
	}
	stop := SimpleRun(2, work)
	// fill jobs and get results at once
	for i := 0; i < TOTALBENCH; i++ {
		select {
		case jobs <- i:
		case <-results:
			i++
		}
	}
	close(stop)
}

func BenchmarkSimpleRun_with_results(*testing.B) {
	ExampleSimpleRun_singleloop()
}

func BenchmarkSimpleRun(*testing.B) {
	var jobs = make(chan int)
	work := func(stop <-chan bool) {
		for {
			select {
			case <-jobs:
			case <-stop:
				return
			}
		}
	}
	stop := SimpleRun(2, work)
	for i := 0; i < TOTALBENCH; i++ {
		select {
		case jobs <- i:
			i++
		}
	}
	close(stop)
}
