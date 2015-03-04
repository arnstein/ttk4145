package statemachine

import (
    "fmt"
)

currentState = INITIALIZE;

/*
States:
INITIALIZE ?
IDLE
MOVING
DOOROPEN
*/
func StateMachine() {

    select {
    case <-signalChannel:
        switch stateMachine {
        case INITIALIZE:
            fmt.Println("Initializing")
            if signalChannel == "ready":
                readyToGetOrders()
        case IDLE:
            fmt.Println("Idle")
        case MOVING:
            fmt.Println("Moving")
        case DOOROPEN:
            fmt.Println("Door open")
        }
    }
}
