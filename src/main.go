package main

import (
	"iohandler"
	"fmt"
	"network"
    "statemachine"
)

func main() {
	arrived := make(chan int)
	orders := make(chan int)

	controlloop.InitCtrl(arrived, orders)
	fmt.Println("int done")
	for {

		orders <- randnr
		<-arrived
	}
}
