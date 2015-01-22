package main

import (
	"fmt"
	"runtime"
	"sync"
)

var mutex = &sync.Mutex{}
var i int = 0




func countUp(done chan bool) {
	for k := 0; k <= 1000000; k++ {
		mutex.Lock()
		i += 1
		mutex.Unlock()
	}
	done<- true
}

func countDown(done chan bool) {
	for k := 0; k <= 1000001; k++ {
		mutex.Lock()
        	i -= 1
        	mutex.Unlock()
	}
	done<- true
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	done1 := make(chan bool, 1)
	done2 := make(chan bool, 1)
	
	go countUp(done1)
	go countDown(done2)
	<-done1
	<-done2
	fmt.Println(i)
}
