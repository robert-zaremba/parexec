package parexec

import (
	"sync"
)

// SimpleWorker is a job function type which is responsible for performing a
// unit of work. It takes a `stop` channel as an agrument. The function
// should listen on it. When stop channel is closed (it returns something) then
// the function should return as well. Check `SimpleRun` function for an example.
type SimpleWorker func(<-chan bool)

// SimpleRun runs `num` workers in separate go routines. You must provide
// a `done`` function which will close all resources. `done` will be called
// when all workers will finish.
// It returns a function to stop all running workers
// Usually you provide a clojure as a worker. Check worker_test.go for examples.
func SimpleRun(num int, w SimpleWorker, done func()) chan bool {
	stop := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			w(stop)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		done()
	}()
	return stop
}
