package parexec

import (
	"strconv"
	"sync"
	"time"
)

type semaphore struct {
	max      int
	acquired int
	cond     *sync.Cond
}

// NewSemaphore returns new Semaphore, which implements a semaphore pattern.
// It is used to limit the number of used resources to `max`. To ask for a resource
// you firstly need to call `Asquire` and when you finish a job you need to release
// a right for resource by calling `Release`.
func NewSemaphore(max int) *semaphore {
	if max <= 0 {
		panic("number of available resources must be positive")
	}
	return &semaphore{
		max:  max,
		cond: sync.NewCond(new(sync.Mutex)),
	}
}

// Acquire locks 1 resource. Blocks if there not enought free resources
func (s *semaphore) Acquire() {
	s.AcquireN(1)
}

// AcquireN locks `n` resources. Blocks if there not enought free resources
func (s *semaphore) AcquireN(n int) {
	if n > s.max {
		panic("can't acquire more resources then `max`=" + strconv.Itoa(s.max))
	}
	s.cond.L.Lock()
	s.acquired += n
	for s.acquired > s.max {
		s.cond.Wait()
	}
	s.cond.L.Unlock()
}

// AcquireCancel is like AcquireN with option to send a signal to cancel the
// request. Returns true if successful and false if cancel occurs.
// Check `AcquireTimeout` implementation for an example use.
func (s *semaphore) AcquireCancel(n int, cancel <-chan struct{}) bool {
	var ok = make(chan struct{})
	var notCanceled = true
	go func() {
		s.cond.L.Lock()
		for s.acquired+n > s.max && notCanceled {
			s.cond.Wait()
		}
		if notCanceled {
			s.acquired += n
		}
		s.cond.L.Unlock()
		ok <- struct{}{}
	}()
	select {
	case <-ok:
	case <-cancel:
		notCanceled = false
		s.cond.Signal()
	}
	return notCanceled
}

// AcquireTimeout is a a usefull pattern to acquire a resource with timetout
// after `d` duration. If duration is super small (ns magnitude) then there is
// a risk that a code will not have a chance to ask for resource befor timeout.
// It's implemented using `AcquireCancel`
func (s *semaphore) AcquireTimeout(n int, d time.Duration) bool {
	timeout := make(chan struct{})
	go func() {
		time.Sleep(d)
		timeout <- struct{}{}
	}()
	return s.AcquireCancel(n, timeout)
}

// Release unlocks 1 resources.
func (s *semaphore) Release() {
	s.ReleaseN(1)
}

// ReleaseN unlocks `n` resources.
func (s *semaphore) ReleaseN(n int) {
	s.cond.L.Lock()
	s.acquired -= n
	if s.acquired <= 0 {
		s.acquired = 0
	}
	if s.acquired <= s.max {
		s.cond.Signal()
	}
	s.cond.L.Unlock()
}

// Wait waits till all resaurces will be available
func (s *semaphore) Wait() {
	s.cond.L.Lock()
	for s.acquired != 0 {
		s.cond.Wait()
	}
	s.cond.L.Unlock()
}

// Available returns number of available resources
func (s *semaphore) Available() int {
	n := s.max - s.acquired
	if n < 0 {
		return 0
	}
	return n
}
