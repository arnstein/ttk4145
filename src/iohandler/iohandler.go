package iohandler

import (
	"driver"
	"globals"
	"network"
	"queue"
	"time"
)

func PollButtons() {

	var buttonStates [3][globals.NUM_FLOORS]int

	for {
		for i := 0; i < globals.NUM_FLOORS; i++ {

			newState := driver.GetButtonSignal(driver.BUTTON_COMMAND, i)

			if newState == 1 && buttonStates[driver.BUTTON_COMMAND][i] == 0 {
				queue.UpdateInsideOrder(i, 1)
			}
			buttonStates[driver.BUTTON_COMMAND][i] = newState
		}
		for i := 0; i < globals.NUM_FLOORS-1; i++ {

			newState := driver.GetButtonSignal(driver.BUTTON_CALL_UP, i)

			if newState == 1 && buttonStates[driver.BUTTON_CALL_UP][i] == 0 {
				network.NewRequest(i, globals.UP)
			}
			buttonStates[driver.BUTTON_CALL_UP][i] = newState
		}
		for i := 1; i < globals.NUM_FLOORS; i++ {

			newState := driver.GetButtonSignal(driver.BUTTON_CALL_DOWN, i)

			if newState == 1 && buttonStates[driver.BUTTON_CALL_DOWN][i] == 0 {
				network.NewRequest(i, globals.DOWN)
			}
			buttonStates[driver.BUTTON_CALL_DOWN][i] = newState
		}

		time.Sleep(10 * time.Millisecond)
	}

}

func CheckFloor() {
	lastFloor := -1
	for {
		floor := driver.GetFloorSensorSignal()

		if lastFloor == -1 && floor != -1 {
			driver.SetFloorIndicator(floor)
			queue.SetCurrentFloor(floor)
			globals.SignalChannel <- globals.FLOORREACHED
		}
		lastFloor = floor
		time.Sleep(10 * time.Millisecond)
	}
}

func Motor(command int) {
	driver.SetMotorDir(command)
}

// queue stuff
