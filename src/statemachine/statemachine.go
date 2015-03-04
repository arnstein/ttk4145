package statemachine

/*
import (
    "fmt"
)
<<<<<<< Updated upstream
const (
    INITIALIZE
    IDLE
    DOOROPEN
    MOVING
)
/*
init
moveOrder
timerout
emptyqueue
floorReached
*/

currentstate = INITIALIZE


func initialize(signal string) {
    switch signals {
    case "init":
        network.NetworkInit()
        controlloop.InitCtrl(arrived, orders)
    case "floorReached":
        // discuss and define behaviour here, should it go to floor 1 all the time? etc etc
        currentstate = IDLE
    }
}

func idle(signal string) {
    switch signals {
    case "moveOrder":
        currentstate = MOVING
}

func doorOpen(signal string) {
    switch signals {
    case "timerOut":
        currentstate = IDLE
        // register done order so the queueChannel will send next order
    }

}

func moving(signal string) {
    switch signals {
    case "floorReached":
        // if right floor: motorControl(stop), currentState = doorOpen

    }

}


func StateMachine() {

    select {
    case <-signalChannel:
        signal := signalChannel
        switch stateMachine {
        case INITIALIZE:
            initialize(signal)
            }
        case IDLE:
            idle(signal)
        case DOOROPEN:
            doorOpen(signal)
        case MOVING:
            moving(signal)

        }
    }
}
