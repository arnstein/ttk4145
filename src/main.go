package main

import (
	"iohandler"
	"fmt"
	"network"
<<<<<<< HEAD
	//	"time"
=======
    "statemachine"
>>>>>>> origin/arnstein
)

func main() {
	arrived := make(chan int)
	orders := make(chan int)
	controlloop.InitCtrl(arrived, orders)
	fmt.Println("int done")
	for {

		orders <- 1
		<-arrived
	}
}
