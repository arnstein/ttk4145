package iohandler

import (
	"driver"
	//"fmt"
	"time"
)

var nextDestination int = 0

func motorControl(arrived chan<- int) {

	lastPosition := 0

	currentPosition := 0

	for {
		currentPosition = driver.GetFloorSensorSignal()

		if currentPosition > -1 {
			lastPosition = currentPosition
			driver.SetFloorLight(currentPosition)
		}

		if currentPosition == nextDestination {
			time.Sleep(200 * time.Millisecond)

			driver.SetMotorDir(driver.DIR_STOP)

			driver.SetDoorOpenLight(1)
			arrived <- currentPosition
			time.Sleep(2 * time.Second)

		} else {
			driver.SetDoorOpenLight(0)

			if nextDestination > lastPosition {
				driver.SetMotorDir(driver.DIR_UP)
			}

			if lastPosition > nextDestination {
				driver.SetMotorDir(driver.DIR_DOWN)
			}
		}
	}
}

func listenForOrders(order <-chan int) {

	for {
		nextDestination = <-order
	}
}

func InitCtrl(arrived chan<- int, orders <-chan int) {

	driver.ElevInit()

	go motorControl(arrived)
	go listenForOrders(orders)

}
