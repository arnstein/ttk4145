package statemachine

import (
	"fmt"

	"driver"
	"globals"
	"iohandler"
	"network"
	"queue"
	"time"
)

/*
init
moveOrder
timerout
emptyqueue
floorReached
*/

const (
	INITIALIZE = 0
	IDLE       = 1
	DOOROPEN   = 2
	MOVING     = 3
)

var currentState int = INITIALIZE

func initialize(signal int) {
	switch signal {
	case globals.INIT:
		network.NetworkInit()
		driver.ElevInit()
		network.InitializeCostsOfOrders() // move this to queue
		go iohandler.PollButtons()
		go iohandler.CheckFloor()
		currentState = IDLE
		fmt.Println("init done")
		globals.SignalChannel <- globals.CHECKORDER
	}
}

func idle(signal int) {
	switch signal {
	case globals.CHECKORDER:
		fmt.Println("got checkorder")
		floor := queue.GetNextOrder()
		if floor != -1 {
			globals.SignalChannel <- globals.MOVEORDER
		} else {
			fmt.Println("ququee  empty")
		}

	case globals.MOVEORDER:
		fmt.Println("got moveorder")
		if queue.GetDirection() == 0 {
			fmt.Println("arrived!")
			currentState = DOOROPEN
			globals.SignalChannel <- globals.FLOORREACHED
		} else {
			fmt.Println("start moving!")
			iohandler.Motor(queue.GetDirection())
			currentState = MOVING
		}
	}
}

func doorOpen(signal int) {
	switch signal {
	case globals.FLOORREACHED:

		//send message that order is handled
		driver.SetDoorOpenLight(1)

		time.Sleep(1 * time.Second)
		floor, direction := queue.CurrentIndexToFloorAndDirection()
		network.RequestServed(floor, direction)

		globals.SignalChannel <- globals.TIMEROUT
	case globals.TIMEROUT:
		driver.SetDoorOpenLight(0)
		currentState = IDLE
		globals.SignalChannel <- globals.CHECKORDER
	}
}

func moving(signal int) {
	switch signal {
	case globals.FLOORREACHED:
		fmt.Println("Floor reached")
		if queue.GetDirection() == 0 {
			iohandler.Motor(globals.STOP)
			currentState = DOOROPEN
			globals.SignalChannel <- globals.FLOORREACHED
		}
	}
}

func StateMachine() {
	for {

		fmt.Println("waiting for signaø")
		signal := <-globals.SignalChannel
		switch currentState {
		case INITIALIZE:
			initialize(signal)
		case IDLE:
			idle(signal)
		case DOOROPEN:
			doorOpen(signal)
		case MOVING:
			moving(signal)
		}
	}
}
