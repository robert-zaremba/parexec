package parexec

import (
	"fmt"
	"sync"
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
	done := func() {
		close(jobs)
	}
	stop := SimpleRun(2, work, done)
	// fill jobs
	for i := 0; i < total; i++ {
		jobs <- &Job{Name: fmt.Sprint(i), Duration: i * 1000}
	}
	// you can't use this: close(results); for r := range results
	// because close(results) is not safe - there might be active goroutine
	// which didn't get stop signal and tries to write to results
	for i := 0; i < total; i++ {
		r := <-results
		fmt.Println(r)
	}
	close(stop)
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
	var results = make(chan int)
	work := func(stop <-chan bool) {
		for {
			select {
			case j := <-jobs:
				results <- j
			case <-stop:
				return
			}
		}
	}
	done := func() {
		close(results)
	}
	stop := SimpleRun(2, work, done)
	sum, r := 0, 0
	for i := 0; i < TOTALBENCH; {
		select {
		case jobs <- i:
			i++
		case r = <-results:
			sum += r
		}
	}
	close(stop) // we need to manually stop to automatically stop results when all workers will finish
	for r = range results {
		sum += r
	}
}

// ExampleSimpleRun_raw present the same computation as ExampleSimpleRun_singleloop
// without using SimpleRun function and stop channel.
func ExampleSimpleRun_raw() {
	var jobs = make(chan int)
	var results = make(chan int)
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			for j := range jobs {
				results <- j
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	sum, r := 0, 0
	for i := 0; i < TOTALBENCH; {
		select {
		case jobs <- i:
			i++
		case r = <-results:
			sum += r
		}
	}
	close(jobs) // we need to close jobs to stop workers and close results
	for r = range results {
		sum += r
	}
}

func BenchmarkSimpleRun(*testing.B) {
	ExampleSimpleRun_singleloop()
}

func BenchmarkSimpleRun_raw(*testing.B) {
	ExampleSimpleRun_raw()
}
