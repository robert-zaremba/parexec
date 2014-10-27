package parexec

import (
	"strconv"
	"sync"
)

type Semaphore struct {
	max      int
	acquired int
	cond     *sync.Cond
}

func NewSemaphore(max int) *Semaphore {
	if max <= 0 {
		panic("number of available resources must be positive")
	}
	return &Semaphore{
		max:  max,
		cond: sync.NewCond(new(sync.Mutex)),
	}
}

// Acquire locks 1 resource. Blocks if there not enought free resources
func (s *Semaphore) Acquire() {
	s.AcquireN(1)
}

// AcquireN locks `n` resources. Blocks if there not enought free resources
func (s *Semaphore) AcquireN(n int) {
	if n > s.max {
		panic("can't acquire more resources then `max`=" + strconv.Itoa(s.max))
	}
	s.cond.L.Lock()
	s.acquired += n
	for s.acquired >= s.max {
		s.cond.Wait()
	}
	s.cond.L.Unlock()
}

// AcquireCancel is like AcquireN with option to send a signal to cancel the
// request. Returns true if successful and false if cancel occurs.
// One useful use-case is a timeout:
//     timeout := make(chan bool)
//     go func() {
//         time.Sleep(d)
//         timeout <- true
//     }()
//     s.AcquireCancel(2, timeout)
func (s *Semaphore) AcquireCancel(n int, cancel <-chan struct{}) bool {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	var ok = make(chan struct{})
	var canceled = false
	go func() {
		for s.acquired+n >= s.max && !canceled {
			s.cond.Wait()
		}
		ok <- struct{}{}
	}()
	select {
	case <-ok:
		s.acquired += n
		return true
	case <-cancel:
		canceled = true
		return false
	}
}

// Release unlocks 1 resources.
func (s *Semaphore) Release() {
	s.ReleaseN(1)
}

// ReleaseN unlocks `n` resources.
func (s *Semaphore) ReleaseN(n int) {
	s.cond.L.Lock()
	s.acquired -= n
	if s.acquired <= 0 {
		s.cond.Signal()
		s.acquired = 0
	}
	s.cond.L.Unlock()
}

// Wait waits till all resaurces will be available
func (s *Semaphore) Wait() {
	s.cond.L.Lock()
	for s.acquired != 0 {
		s.cond.Wait()
	}
	s.cond.L.Unlock()
}

// Available returns number of available resources
func (s *Semaphore) Available() int {
	n := s.max - s.acquired
	if n < 0 {
		return 0
	}
	return n
}
