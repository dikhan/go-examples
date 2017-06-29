package main

import (
	"errors"
	"fmt"
	"time"
)

type Promise struct {
	successChannel chan interface{}
	failureChannel chan error
}

type PurshaseOrder1 struct {
	Number int
	Value  float64
}

func main() {
	po := new(PurshaseOrder1)
	po.Value = 42.27

	SavePO1(po, false).
		Then(
			func(obj interface{}) error {
				po := obj.(*PurshaseOrder1)
				fmt.Printf("Purchase Order saved with ID: %d\n", po.Number)
				return nil
			},
			func(err error) {
				fmt.Printf("Failed to save Purchase Order: " + err.Error() + "\n")
			}).
		Then(
			func(obj interface{}) error {
				fmt.Println("Second promise success")
				return nil
			},
			func(err error) {
				fmt.Println("Second promise failed: " + err.Error())
			})
	fmt.Scanln()
}

func SavePO1(po *PurshaseOrder1, shouldFail bool) *Promise {
	result := new(Promise)
	// Creating the buffer allows the consumer of this method to decide whether they want to chain a new function or not. If
	// they don't then these channels will never be drained and they would block the applications if they were unbuffered.
	result.successChannel = make(chan interface{})
	result.failureChannel = make(chan error)

	// This time will simulate the asynchronous behaviour by creating a go routine
	go func() {
		time.Sleep(2 * time.Second)
		if shouldFail {
			result.failureChannel <- errors.New("Fail to save purchase order")
		} else {
			po.Number = 1234
			result.successChannel <- po
		}
	}()
	return result
}

// Then method accepting two functions as parameters which will be called when the promise is finished. The method
// also returns another promise to allow the chaining of more promises afterwards
func (p *Promise) Then(success func(interface{}) error, failure func(error)) *Promise {

	// Create the promise that the method is going to return
	result := new(Promise)

	// By giving a small buffer size we can bu sure that the method doesnt stop executing if another handler is not chained to this one
	result.successChannel = make(chan interface{}, 1)
	result.failureChannel = make(chan error, 1)

	timeout := time.After(1 * time.Second)

	// Now we add the code that decides if the promise was successful or not
	// This will allow us to return the new promise synchronously while it is actually being process asynchronosly
	go func() {
		select {
		case obj := <-p.successChannel:
			newErr := success(obj)
			if newErr == nil {
				result.successChannel <- obj
			} else {
				result.failureChannel <- newErr
			}
		case err := <-p.failureChannel:
			failure(err)
			result.failureChannel <- err
		case <-timeout:
			failure(errors.New("Promise timeout"))
		}
	}()
	return result
}
