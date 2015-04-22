package queue

import (
	"driver"
	"fmt"
	"globals"
	"io/ioutil"
	"os"
	"time"
)

const (
	FLOORS            = 4
	ORDERS_ARRAY_SIZE = (2 * FLOORS) - 2

	GLOBAL = 1
	LOCAL  = 2
	NONE   = 0

	DIR_DOWN = -1
	DIR_ANY  = 0
	DIR_UP   = 1
)

var orders [ORDERS_ARRAY_SIZE]int

var insideOrders [FLOORS]byte

var OrderBackup [ORDERS_ARRAY_SIZE]time.Time

var position int
var currentFloor int

func UpdateInsideOrder(floor int, status int) {

	if byte(status) == insideOrders[floor] {
		return
	}

	insideOrders[floor] = byte(status)
	ioutil.WriteFile("insideQueue", insideOrders[:], 0666)
	// should the queue really use driver?
	driver.SetButtonLamp(driver.BUTTON_COMMAND, floor, status)
	globals.SignalChannel <- globals.CHECKORDER

}

func Init() {
	array, err := ioutil.ReadFile("insideQueue")
	if err != nil {
		_, err := os.Create("insideQueue")
		globals.CheckError(err)
	}
	globals.CheckError(err)
	if len(array) != FLOORS {
		for i := 0; i < FLOORS; i++ {
			UpdateInsideOrder(i, 0)
		}
	} else {
		for i := 0; i < FLOORS; i++ {
			status := int(array[i])
			UpdateInsideOrder(i, status)
		}
	}
}

func IndexToFloorAndDirection(index int) (int, int) {

	if index == -1 {
		index = position
	}

	if index < ORDERS_ARRAY_SIZE/2 {
		return index, 1

	}
	return ORDERS_ARRAY_SIZE - index, -1

}

func PrintQueue() {

	for i := 0; i < ORDERS_ARRAY_SIZE; i++ {
		fmt.Print(" ")
		fmt.Print(orders[i])
	}

	fmt.Println()

	for i := 0; i < ORDERS_ARRAY_SIZE; i++ {
		if position == i {
			fmt.Print("-!")
		} else {
			fmt.Print("--")
		}
	}
	fmt.Println()

}

func SetCurrentFloor(floor int) {
	currentFloor = floor

	// move to the next request
	for i := 0; i < ORDERS_ARRAY_SIZE; i++ {
		position = (position + 1) % ORDERS_ARRAY_SIZE
		nextFloor, _ := IndexToFloorAndDirection(position)
		if orders[position] != 0 || insideOrders[nextFloor] == 1 {
			break
		}
	}
	// move back to appearance of the floor
	for i := 0; i < ORDERS_ARRAY_SIZE; i++ {
		nextFloor, _ := IndexToFloorAndDirection(position)
		if nextFloor == floor {
			return
		}
		// avoid negative modulo
		//if position == floor || position == ORDERS_ARRAY_SIZE-floor {
		//return
		//}
		position = (position - 1 + ORDERS_ARRAY_SIZE) % ORDERS_ARRAY_SIZE
	}
}

func FloorAndDirToIndex(floor int, dir int) int {

	wayUp := floor
	wayDown := ORDERS_ARRAY_SIZE - floor

	if dir == DIR_UP {
		return wayUp
	}

	if dir == DIR_DOWN {
		return wayDown
	}

	if wayUp < position && position < wayDown {
		return wayDown
	}

	return wayUp
}

func AddToQueue(floor int, dir int, globalOrLocal int) {
	index := FloorAndDirToIndex(floor, dir)

	if orders[index] == GLOBAL {
		return
	}
	orders[index] = globalOrLocal
	globals.SignalChannel <- globals.CHECKORDER
}

func RemoveFromQueue(floor int, dir int) {

	orders[FloorAndDirToIndex(floor, dir)] = NONE

}

func AddToBackupQueue(floor int, dir int) {
	index := FloorAndDirToIndex(floor, dir)

	OrderBackup[index] = time.Now()
}

func RemoveFromBackupQueue(floor int, dir int) {

	index := FloorAndDirToIndex(floor, dir)
	OrderBackup[index] = time.Unix(0, 0)

}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func CalculateCost(floor int, dir int) int {

	endIndex := FloorAndDirToIndex(floor, dir)

	cost := abs(currentFloor - floor)

	for i := position; i != endIndex; i = (i + 1) % ORDERS_ARRAY_SIZE {
		if orders[i] != NONE {
			cost++
		}
	}
	return cost
}

func GetNextOrder() int {

	nextOrder := -1

	//get index of next order
	for i := 0; i < ORDERS_ARRAY_SIZE; i++ {

		index := (i + position) % ORDERS_ARRAY_SIZE
		floor, _ := IndexToFloorAndDirection(index)

		if orders[index] != NONE || insideOrders[floor] == byte(1) {
			nextOrder = index
			break
		}
	}

	if nextOrder == -1 {
		return -1
	}

	nextFloor, _ := IndexToFloorAndDirection(nextOrder)

	// detect special case: up and down order for current floor
	if currentFloor == nextFloor && nextOrder != position {
		return -2
	}

	return nextFloor
}

func GetDirection() int {

	next := GetNextOrder()
	dir := next - currentFloor
	if dir > 0 {
		dir = 1
	}
	if dir < 0 {
		dir = -1
	}

	if dir != 0 {
		floor, _ := IndexToFloorAndDirection(-1)
		position = FloorAndDirToIndex(floor, dir)
	}

	fmt.Print("we are in floor ")
	fmt.Print(currentFloor)
	fmt.Print(" and index ")
	fmt.Print(position)
	fmt.Print(" and want to go to ")
	fmt.Println(next)

	return dir
}
