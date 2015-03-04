package iohandler

import (
	"driver"
	//"fmt"
	"time"
)

var nextDst int = 0

func motorCtrl(arrived chan<- int) {

	lastPos := 0

	currentPos := 0

	for {
		currentPos = driver.GetFloorSensorSignal()

		if currentPos > -1 {
			lastPos = currentPos
			driver.SetFloorLight(currentPos)
		}

		if currentPos == nextDst {
			time.Sleep(200 * time.Millisecond)

			driver.SetMotorDir(driver.DIR_STOP)

			driver.SetDoorOpenLight(1)
			arrived <- currentPos
			time.Sleep(2 * time.Second)

		} else {
			driver.SetDoorOpenLight(0)

			if nextDst > lastPos {
				driver.SetMotorDir(driver.DIR_UP)
			}

			if lastPos > nextDst {
				driver.SetMotorDir(driver.DIR_DOWN)
			}
		}
	}
}

func listenForOrders(order <-chan int) {

	for {
		nextDst = <-order
	}
}

func InitCtrl(arrived chan<- int, orders <-chan int) {

	driver.ElevInit()

	go motorCtrl(arrived)
	go listenForOrders(orders)

}
