package parexec

import (
	"fmt"
	"testing"
)

func ExampleNewSemaphore() {
	sem := NewSemaphore(1)
	total := 4
	job := func(i int) {
		fmt.Println(i)
		sem.Release()
	}
	for i := 0; i < total; i++ {
		sem.Acquire()
		go job(i)
	}
	sem.Wait()
	// Output:
	// 0
	// 1
	// 2
	// 3
}

func TestAcquireTimeout(t *testing.T) {
	sem := NewSemaphore(1)
	sem.Acquire()
	ok := sem.AcquireTimeout(1, 10000)
	if ok {
		t.Error("Should be timeouted")
	}
	sem.Release()
	ok = sem.AcquireTimeout(1, 10000)
	if !ok {
		t.Error("Shouldn't be timeouted")
	}
	sem.Release()
}

func BenchmarkSemaphore(b *testing.B) {
	sem := NewSemaphore(2)
	job := func(i int) {
		sem.Release()
	}
	for i := 0; i < TOTALBENCH; i++ {
		sem.Acquire()
		go job(i)
	}
}
