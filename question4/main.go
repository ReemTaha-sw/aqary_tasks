package main

import (
	"fmt"
	"sync"
	// "time"
)

const bufferSize = 10

func main() {
	M := 2 // Number of reading goroutines
	N := 8 // Number of writing goroutines

	buffer := make([]byte, bufferSize)
	reading := make(chan struct{}, 1)
	writeing := make(chan struct{}, 1)

	// Initialize read channel to allow the first read
	reading <- struct{}{}

	// Writing goroutine
	for i := 0; i < N; i++ {
		go func(id int) {
			for {
				// Write to buffer
				writeing <- struct{}{}
				buffer[id%bufferSize] = byte(id)
				fmt.Printf("Writing %d wrote: %d\n", id, buffer[id%bufferSize])
				// time.Sleep(time.Millisecond * 500)
				<-reading
			}
		}(i)
	}

	// Reading goroutines
	var waitgroup sync.WaitGroup
	for i := 0; i < M; i++ {
		waitgroup.Add(1)
		go func(id int) {
			defer waitgroup.Done()
			for {
				// Read from buffer
				<-writeing
				value := buffer[id%bufferSize]
				fmt.Printf("Reader %d read: %d\n", id, value)
				// time.Sleep(time.Millisecond * 500)
				reading <- struct{}{}
			}
		}(i)
	}

	// time.Sleep(time.Second * 3)

	// Wait for all reading goroutines to finish
	waitgroup.Wait()
}
