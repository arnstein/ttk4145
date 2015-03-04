package main

import (
	"controlloop"
	"fmt"
	"network"
    "statemachine"
)

func main() {
	network.NetworkInit()
	arrived := make(chan int)
	orders := make(chan int)

	controlloop.InitCtrl(arrived, orders)
	fmt.Println("int done")
	for {

		orders <- randnr
		<-arrived
	}
}
