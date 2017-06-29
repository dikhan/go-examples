package main

import "fmt"

// Using channels as a stand in for the callback function as opposed to passing in another function
// This decouples the two functions and gives flexibility. It also enhances the testability of the first function
// since we can easily send and receive mock data through channels; whereas using a function can make it difficult to
// separate the logic of the two functions
func main() {
	po := new(PurshaseOrder)
	po.Value = 42.27

	ch := make(chan *PurshaseOrder)

	go SavePO(po, ch)

	newPo := <- ch
	fmt.Printf("PO: %v", newPo)

	// This way the callback (channel) can be shared with multiple receivers and the number
	// of go routines can increase - allowing a faster draining of the channel and balance the load
	// of the application
}

type PurshaseOrder struct {
	Number int
	Value float64
}

func SavePO(po *PurshaseOrder, callback chan *PurshaseOrder) {
	po.Number = 1234
	callback <- po
}