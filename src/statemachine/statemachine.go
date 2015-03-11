package statemachine

import (
	//"fmt"
	//"iohandler"
	"network"
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
	case INIT:
		network.NetworkInit()
        driver.ElevInit()
        go iohandler.PollButtons()
        // all the other inits
		currentState = IDLE
	}
}

func idle(signal int) {
	switch signal {
	case MOVEORDER:
        direction = queue.getDirection()
        motor(direction)
		currentState = MOVING
	}
}

func doorOpen(signal int) {
	switch signal {
	case TIMEROUT:
        queue.removeFromQueue()
		currentState = IDLE
	}
}

func moving(signal int) {
	switch signal {
	case FLOORREACHED:
        if queue.rightFloor() == 1 {
            motor(STOP)
            currentState = DOOROPEN
	    }
    }
}

func StateMachine(signalChannel <-chan int) {
    select {
    case signal := <-signalChannel:
		switch signal {
		case globals.INITIALIZE:
			initialize(signal)
		case globals.IDLE:
			idle(signal)
		case globals.DOOROPEN:
			doorOpen(signal)
		case globals.MOVING:
			moving(signal)

		}
	}
}

