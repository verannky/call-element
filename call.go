package main

import (
	"fmt"
	"sync"
)

// Global Variables
var x int
var y int
var xData []int
var yData []int
var eventData []int

// Channel
type variable struct {
	event int
	done  chan bool // Added a channel to signal completion
}

// Call element
func CallElement(c chan variable) {
	for {
		data := <-c

		if data.event == 0 {
			// Handle event 0
			if x == 0 {
				x = 1
			} else {
				x = 0
			}
		} else if data.event == 1 {
			// Handle event 1
			if y == 0 {
				y = 1
			} else {
				y = 0
			}
		}

		// Record data for history
		xData = append(xData, x)
		yData = append(yData, y)
		eventData = append(eventData, data.event)

		// Signal completion
		data.done <- true
	}
}

// main function
func main() {
	c := make(chan variable)
	var wg sync.WaitGroup

	go CallElement(c)

	events := []int{1, 1, 0, 0, 1, 1, 0, 1}

	for _, event := range events {
		wg.Add(1)
		done := make(chan bool)

		msg := variable{event: event, done: done}
		c <- msg

		<-done
		wg.Done()
	}

	// Wait for all events to complete
	wg.Wait()

	// Close Application
	fmt.Println("Press the Enter Key to stop anytime")
	fmt.Scanln()

	// Print History
	fmt.Println("Input:", eventData)
	fmt.Println("Kiri  :", xData)
	fmt.Println("Kanan :", yData)
}
