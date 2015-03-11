package iohandler

import (
	"driver"
	"fmt"
	//"time"
	"globals"
	"queue"
)

func PollButtons() {

	var buttonStates [3][globals.NUM_FLOORS]int

	for {
		for i := 0; i < globals.NUM_FLOORS; i++ {

			newState := driver.GetButtonSignal(driver.BUTTON_COMMAND, i)

			if newState == 1 && buttonStates[driver.BUTTON_COMMAND][i] == 0 {
				fmt.Println("order local " + string(i))
			}
			buttonStates[driver.BUTTON_COMMAND][i] = newState
		}
		for i := 0; i < globals.NUM_FLOORS-1; i++ {

			newState := driver.GetButtonSignal(driver.BUTTON_CALL_UP, i)

			if newState == 1 && buttonStates[driver.BUTTON_CALL_UP][i] == 0 {
				fmt.Println("order up    " + string(i))
			}
			buttonStates[driver.BUTTON_CALL_UP][i] = newState
		}
		for i := 1; i < globals.NUM_FLOORS; i++ {

			newState := driver.GetButtonSignal(driver.BUTTON_CALL_DOWN, i)

			if newState == 1 && buttonStates[driver.BUTTON_CALL_DOWN][i] == 0 {
				fmt.Println("order down  " + string(i))
			}
			buttonStates[driver.BUTTON_CALL_DOWN][i] = newState
		}

	}
}

func checkFloor(signalChannel chan<- int) {
	lastFloor := -1
	for {
		floor := driver.GetFloorSensorSignal()

		if lastFloor == -1 && floor != -1 {
            driver.SetFloorIndicator(floor)
			queue.SetCurrentFloor(floor)
			signalChannel <- globals.FLOORREACHED
		}
		lastFloor = floor
	}
}

func motor(command int) {
    SetMotorDir(command)
}

// queue stuff
