package iohandler

import (
	"driver"
	"globals"
	"queue"
	"time"
)

func IoInit() {
	driver.ElevInit()
	go SetLights()
}

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
				globals.NewRequest <- [2]int{i, globals.UP}

			}
			buttonStates[driver.BUTTON_CALL_UP][i] = newState
		}
		for i := 1; i < globals.NUM_FLOORS; i++ {

			newState := driver.GetButtonSignal(driver.BUTTON_CALL_DOWN, i)

			if newState == 1 && buttonStates[driver.BUTTON_CALL_DOWN][i] == 0 {
				globals.NewRequest <- [2]int{i, globals.DOWN}
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

func SetLights() {

	var lightData [3]int
	for {
		lightData = <-globals.LightsChannel

		lightType := lightData[0]
		floor := lightData[1]
		value := lightData[2]

		switch lightType {

		case driver.DIR_UP:
			driver.SetButtonLamp(driver.BUTTON_CALL_UP, floor, value)

		case driver.DIR_DOWN:
			driver.SetButtonLamp(driver.BUTTON_CALL_DOWN, floor, value)

		case driver.BUTTON_COMMAND:
			driver.SetButtonLamp(driver.BUTTON_COMMAND, floor, value)

		case driver.DOOROPEN:
			driver.SetDoorOpenLight(value)
		}
	}
}
