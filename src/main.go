package main

import (
	"controlloop"
	"fmt"
	"network"
	//	"time"
)

func main() {
	network.NetworkInit()
	arrived := make(chan int)
	orders := make(chan int)
	controlloop.InitCtrl(arrived, orders)
	fmt.Println("int done")
	for {

		orders <- 1
		<-arrived
	}
}
