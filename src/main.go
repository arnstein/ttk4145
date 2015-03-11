package main

import (
	//"driver"
	"fmt"
	//"iohandler"
	//"network"
	"queue"
	//"statemachine"
	//"network/udp"
	//	"time"
	//"statemachine"
	//	"net"
)

func main() {
	//	go statemachine.StateMachine()
	//driver.ElevInit()
	//network.NetworkInit()
	//go iohandler.PollButtons()
	//arrived := make(chan int)
	//orders := make(chan int)
	//fmt.Println("int done")
	//for {
	//orders <- 1
	//<-arrived
	//}
	fmt.Println(" 0 1 2 3 2 1")
	fmt.Println()
	queue.PrintQueue()

	queue.AddToQueue(1, 1, 1)
	queue.PrintQueue()

	queue.AddToQueue(1, 0, 2)
	queue.PrintQueue()

	queue.AddToQueue(2, -1, 2)
	queue.PrintQueue()

	queue.SetCurrentFloor(0)
	queue.PrintQueue()

	queue.SetCurrentFloor(1)
	queue.PrintQueue()
	queue.SetCurrentFloor(2)
	queue.PrintQueue()
	fmt.Println(queue.GetNextOrder())
	fmt.Println(queue.CalculateCost(1, -1))
}
