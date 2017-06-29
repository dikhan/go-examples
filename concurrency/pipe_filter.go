package main

import "fmt"

func main() {
	ch := make(chan int)
	go generate(ch)
	for {
		prime := <- ch
		fmt.Println(prime)
		out := make(chan int)
		go filter(ch, out, prime)
		ch = out
	}
}

func generate(ch chan int) {
	for i := 2; i<10 ; i++ {
		ch <- i
	}
}

func filter(in, out chan int, prime int) {
	for {
		i := <-in
		if i % prime != 0 {
			fmt.Printf("Filter by %d - new prime %d\n", prime, i)
			out <- i
		} else {
			fmt.Printf("Discarted - %d\n", i)
		}
	}
}
