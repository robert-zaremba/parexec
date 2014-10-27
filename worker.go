package parexec

// SimpleWorker is a job function type which is responsible for performing a
// unit of work. It takes a `stop` channel as an agrument. The function
// should listen on it. When stop channel is closed (it returns something) then
// the function should return as well. Check `SimpleRun` function for an example.
type SimpleWorker func(<-chan bool)

// SimpleRun runs `num` workers in separate go routines.
// it returns a function to stop all running workers
// Usually you provide a clojure as a worker. Check worker_test.go for examples.
func SimpleRun(num int, w SimpleWorker) chan bool {
	done := make(chan bool)
	for i := 0; i < num; i++ {
		go w(done)
	}
	return done
}
