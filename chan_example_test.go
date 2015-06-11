package parexec

import (
	"fmt"
	"time"
)

func do(c chan int, name string) {
	for {
		x, ok := <-c
		fmt.Println(name, x, ok)
		time.Sleep(2 * time.Millisecond)
	}
}

func ExampleChan() {
	c := make(chan int)
	go do(c, "A")
	time.Sleep(2 * time.Millisecond)
	c <- 2
	c <- 3
	// after close, `<-c` returns (zero value, false) all the time immediately.
	close(c)
	// you can't send to closed channel. This panics!
	// c <- 4
	time.Sleep(4 * time.Millisecond)
	println("--  resetting channel to nil")
	c = nil
	go do(c, "B") // c is a new channel, let's consume it!
	go func() {
		c <- 3
		println("we can send to nil channel, but it will block!")
	}()
	go func() {
		<-c
		println("receiving from nil channel also blocks")
	}()
	time.Sleep(4 * time.Millisecond)

	// Output:
	// A 2 true
	// A 3 true
	// A 0 false
	// A 0 false
	// A 0 false
}
