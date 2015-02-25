package main

import (
	. "fmt"
	"runtime"
	//"time"
)

var count int

func add(ownChan chan int, sync chan<- int) {
	for i := 0; i < 1000*1000; i++ {
		<-ownChan
		count++
		ownChan <- 0
	}
	sync <- 0
}

func sub(ownChan chan int, sync chan<- int) {
	for i := 0; i < 1000*1000-1; i++ {
		<-ownChan
		count--
		ownChan <- 0
	}
	sync <- 0
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	count = 0

	ownChan := make(chan int, 1)
	sync1 := make(chan int)
	sync2 := make(chan int)

	go add(ownChan, sync1)
	go sub(ownChan, sync2)

	ownChan <- 0

	<-sync1
	<-sync2

	Println(count)

}
