package statemachine

import (
	//"fmt"
	//"iohandler"
	"network"
)

const (
	INITIALIZE = 0
	IDLE       = 1
	DOOROPEN   = 2
	MOVING     = 3

	INIT         = 0
	MOVEORDER    = 1
	TIMEROUT     = 2
	EMPTYQUEUE   = 3
	FLOORREACHED = 4
)

/*
init
moveOrder
timerout
emptyqueue
floorReached
*/

var currentstate int = INITIALIZE

func initialize(signal int) {
	switch signal {
	case INIT:
		network.NetworkInit()
		//eventhandler.InitCtrl(arrived, orders)
	case FLOORREACHED:
		// discuss and define behaviour here, should it go to floor 1 all the time? etc etc
		currentstate = IDLE
	}
}

func idle(signal int) {
	switch signal {
	case MOVEORDER:
		currentstate = MOVING
	}
}

func doorOpen(signal int) {
	switch signal {
	case TIMEROUT:
		currentstate = IDLE
		// register done order so the queueChannel will send next order
	}

}

func moving(signal int) {
	switch signal {
	case FLOORREACHED:
		// if right floor: motorControl(stop), currentState = doorOpen

	}

}

func StateMachine(signalChannel <-chan int) {

	for {
		signal := <-signalChannel
		switch signal {
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
