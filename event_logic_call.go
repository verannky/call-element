// Tanggal: 24-10-2023
// Kelompok 2
// 23222305 - Nunky Febilia Verany
// 23222306 - Deddy Welsan
// 23222308 - Muhammad Rizki Putra
// EL5102 Arsitektur Komputer Lanjut
// Topik: Call Element
// Teknik Komputer - Institut Teknologi Bandung

package main

// Here "fmt" is formatted IO which
// is same as C’s printf and scanf.
import (
	"fmt"
	"math/rand"
	"time"
)

type CallInpEvt struct {
	R1 bool
	R2 bool
}

type CallOutEvt struct {
	D1 bool
	D2 bool
}

type SharedEvt struct {
	R bool
	D bool
}

// function that return random number between min and max
func random(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func ReceiveResponse(o chan CallOutEvt) {
	var prevD1, prevD2 bool
	prevD1 = false // assumed start from LOW logic
	prevD2 = false // assumed start from LOW logic
	for {
		channel := <-o
		if channel.D1 != prevD1 {
			prevD1 = channel.D1
			fmt.Printf("Receive response from D1\n\n")
		} else if channel.D2 != prevD2 {
			prevD2 = channel.D2
			fmt.Printf("Receive response from D2\n\n")
		}
	}
}

func SharedResource(s chan SharedEvt) {
	var output SharedEvt
	var prevR, prevD bool
	prevR = false // assumed start from LOW logic
	prevD = false // assumed start from LOW logic
	min := 100    // delay in millisecond
	max := 1000   // delay in millisecond
	for {
		channel := <-s
		fmt.Println("shared resource get data")
		if channel.R != prevR {
			time.Sleep(time.Duration(random(min, max)) * time.Millisecond) // delay from 100 to 1000ms to simulate the process
			prevD = !prevD
			output = SharedEvt{R: prevR, D: prevD}
			prevR = channel.R
			s <- output
			fmt.Println("shared resource return response")
		}
	}
}

func Call(i chan CallInpEvt, o chan CallOutEvt, s chan SharedEvt) {
	var sharedInput SharedEvt
	fromR1 := false
	fromR2 := false
	reqDone := true
	var prevR1, prevR2, prevD1, prevD2, prevR, prevD bool
	prevR1 = false // assumed start from LOW logic
	prevR2 = false // assumed start from LOW logic
	prevD1 = false // assumed start from LOW logic
	prevD2 = false // assumed start from LOW logic
	prevR = false  // assumed start from LOW logic
	prevD = false  // assumed start from LOW logic

	for {
		inpChannel := <-i
		if reqDone {
			if (prevR1 == inpChannel.R1) && (prevR2 == inpChannel.R2) {
				fmt.Printf("not an event\n\n")
			} else {
				randNum := rand.Intn(2)
				if randNum == 0 {
					if prevR1 != inpChannel.R1 {
						fromR1 = true
						fmt.Println("Get request from R1, waiting for response")
					} else if prevR2 != inpChannel.R2 {
						fromR2 = true
						fmt.Println("Get request from R2, waiting for response")
					}
				} else {
					if prevR2 != inpChannel.R2 {
						fromR2 = true
						fmt.Println("Get request from R2, waiting for response")
					} else if prevR1 != inpChannel.R1 {
						fromR1 = true
						fmt.Println("Get request from R1, waiting for response")
					}
				}

				prevR1 = inpChannel.R1
				prevR2 = inpChannel.R2
				prevR = !prevR
				sharedInput = SharedEvt{R: prevR, D: prevD}
				s <- sharedInput
				reqDone = false

				response := <-s
				if response.D != prevD {
					prevD = response.D
					if fromR1 && !fromR2 {
						prevD1 = !prevD1
						output := CallOutEvt{D1: prevD1, D2: prevD2}
						o <- output
						fromR1 = false
						fmt.Println("output response to D1, procedure completed")
					} else if !fromR1 && fromR2 {
						prevD2 = !prevD2
						output := CallOutEvt{D1: prevD1, D2: prevD2}
						o <- output
						fromR2 = false
						fmt.Println("output response to D2, procedure completed")
					}
					reqDone = true
				}
			}
		} else {
			fmt.Println("Previous procedure is not done yet, request is ignored")
		}
	}
}

func main() {

	i := make(chan CallInpEvt)
	o := make(chan CallOutEvt)
	s := make(chan SharedEvt)

	go Call(i, o, s)
	go SharedResource(s)
	go ReceiveResponse(o)

	// Testing
	const data_num = 10
	r1_data := [data_num]bool{false, true, true, false, false, true, true, false, false, true}
	r2_data := [data_num]bool{false, false, true, true, false, false, false, true, true, false}

	for j := 0; j < data_num; j++ {
		x := CallInpEvt{R1: r1_data[j], R2: r2_data[j]}
		time.Sleep(200 * time.Millisecond)
		i <- x
	}

	// waitkey
	time.Sleep(1 * time.Second)
	fmt.Println("Press the Enter Key to stop the program")
	fmt.Scanln()
}
