package parexec

// SimpleWorker is a job function which accepts a *stop* channel and should
// listens on it. When stop channel is closed (it returns something) then
// the function should return.
type SimpleWorker func(chan bool)

// SimpleRun runs `num` workers in separate go routines.
// it returns a function to stop all running workers
// Usually you provide a clojure as a worker:
//
//    jobs    = make(chan Job)
//    results = make(chan bool)
//    func work(done chan bool) {
//    	for {
//    		select {
//    		case j := <-jobs:
//    			results<- doSomething(j)
//    		case <-done:
//    		}
//    	}
//    }
//    stop := SimpleRun(num, work)
//    // fill jobs
//    counter := 0
//    for r := range results {
//    	counter++
//    	if counter > total {
//    		close(stop)
//    		break
//    	}
//    }
func SimpleRun(num int, w SimpleWorker) chan bool {
	done := make(chan bool)
	for i := 0; i < num; i++ {
		go w(done)
	}
	return done
}
