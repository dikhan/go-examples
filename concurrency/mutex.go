package main

import (
	"runtime"
	"fmt"
	"sync"
)

func main() {
	mutex := new(sync.Mutex)
	runtime.GOMAXPROCS(4)
	// Sharing the loop variables with the asynchronous goroutines is not a good idea
	for i:=1; i<10;i++ {
		for j:=1; j<10;j++ {
			mutex.Lock()
			go func() {
				fmt.Printf("%d + %d = %d\n", i, j, i+j)
				mutex.Unlock()
			}()
		}
	}

	// This will force the main thread to stay alive while the other go routines are still working
	fmt.Scanln()
}
