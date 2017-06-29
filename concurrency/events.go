package main

import "fmt"

type Button struct {
	// key: name of the event that is being listened for
	// value: slice of channels that accept the event object
	// The slice will be populated with each object that is listening for the event
	eventListeners map[string][]chan string
}

func main() {
	btn := MakeButton()

	handlerOne := make(chan string)
	handlerTwo := make(chan string)

	btn.AddEvent("click", handlerOne)
	btn.AddEvent("click", handlerTwo)

	go func() {
		for {
			msg := <- handlerOne
			fmt.Println("Handler one: " + msg )
		}
	}()

	go func() {
		for {
			msg := <- handlerTwo
			fmt.Println("Handler two: " + msg )
		}
	}()

	btn.TriggerEvent("click", "button clicked :)")
	btn.RemoveEventListener("click", handlerTwo)
	btn.TriggerEvent("click", "button clicked again!!")

	fmt.Scanln()
}

func MakeButton() *Button {
	result := new(Button)
	result.eventListeners = make(map[string][]chan string)
	return result
}

func (b *Button) AddEvent(event string, responseChannel chan string) {
	if _, present := b.eventListeners[event]; present {
		b.eventListeners[event] = append(b.eventListeners[event], responseChannel)
 	} else {
		b.eventListeners[event] = []chan string{responseChannel}
	}
}

func (b *Button) RemoveEventListener(event string, listenerChannel chan string) {
	if _, present := b.eventListeners[event]; present {
		for idx, _ := range b.eventListeners[event] {
			if b.eventListeners[event][idx] == listenerChannel {
				b.eventListeners[event] = append(b.eventListeners[event][:idx], b.eventListeners[event][idx+1:]...)
				break
			}
		}
	}
}

func (b * Button) TriggerEvent(event string, response string) {
	if _, present := b.eventListeners[event]; present {
		for _, handler := range b.eventListeners[event] {
			go func(handler chan string) {
				handler <- response
			}(handler)
		}
	}
}