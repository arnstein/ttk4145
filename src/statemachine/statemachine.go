package statemachine

import (
	"fmt"
	"globals"
	"iohandler"
	"network"
	"queue"
	"time"
)

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
		fmt.Println("Welcome to the Uppy-Downy Machine(tm)!")
		fmt.Println("Please wait for initialization to be done.")
		network.NetworkInit()
		iohandler.IoInit()
		queue.Init()
		network.InitializeCostsOfOrders()
		go iohandler.PollButtons()
		go iohandler.CheckFloor()
		go network.CheckBackupTimeouts()
		currentState = IDLE
		fmt.Println("Initialization done. Please enjoy the Uppy-Downy Machine(tm).")
		globals.SignalChannel <- globals.CHECKORDER
	}
}

func idle(signal int) {
	switch signal {
	case globals.CHECKORDER:
		floor := queue.GetNextOrder()

		// special case: only up and down in one floor left to be served
		if floor == queue.SPECIAL_CASE_ORDER {
			floor, direction := queue.IndexToFloorAndDirection(-1)
			network.RequestServed(floor, -1*direction)
			currentState = DOOROPEN
			globals.SignalChannel <- globals.FLOORREACHED

		} else if floor >= 0 {
			globals.SignalChannel <- globals.MOVEORDER
		}

	case globals.MOVEORDER:

		if queue.GetDirection() == 0 {
			currentState = DOOROPEN
			globals.SignalChannel <- globals.FLOORREACHED

		} else {
			iohandler.Motor(queue.GetDirection())
			currentState = MOVING
		}
	}
}

func doorOpen(signal int) {
	switch signal {
	case globals.FLOORREACHED:

		floor, direction := queue.IndexToFloorAndDirection(-1)
		queue.UpdateInsideOrder(floor, 0)

		network.RequestServed(floor, direction)
		globals.LightsChannel <- [3]int{3, 0, 1}

		time.Sleep(1 * time.Second)
		globals.LightsChannel <- [3]int{3, 0, 0}

		currentState = IDLE
		globals.SignalChannel <- globals.CHECKORDER
	}
}

func moving(signal int) {

	switch signal {
	case globals.FLOORREACHED:

		if queue.GetDirection() == 0 {
			iohandler.Motor(globals.STOP)
			currentState = DOOROPEN
			globals.SignalChannel <- globals.FLOORREACHED
		}
	}
}

func StateMachine() {
	for {
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
