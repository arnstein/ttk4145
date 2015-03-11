package statemachine

import (
	//"fmt"
	//"iohandler"
	"globals"
	"network"
)

/*
init
moveOrder
timerout
emptyqueue
floorReached
*/

var currentstate int = globals.INITIALIZE

func initialize(signal int) {
	switch signal {
	case globals.INIT:
		network.NetworkInit()
		//eventhandler.InitCtrl(arrived, orders)
	case globals.FLOORREACHED:
		// discuss and define behaviour here, should it go to floor 1 all the time? etc etc
		currentstate = globals.IDLE
	}
}

func idle(signal int) {
	switch signal {
	case globals.MOVEORDER:
		currentstate = globals.MOVING
	}
}

func doorOpen(signal int) {
	switch signal {
	case globals.TIMEROUT:
		currentstate = globals.IDLE
		// register done order so the queueChannel will send next order
	}

}

func moving(signal int) {
	switch signal {
	case globals.FLOORREACHED:
		// if right floor: motorControl(stop), currentState = doorOpen

	}

}

func StateMachine(signalChannel <-chan int) {

	for {
		signal := <-signalChannel
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
