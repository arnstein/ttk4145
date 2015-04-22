package queue

import (
	"fmt"
	"globals"
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

var OrderBackup [ORDERS_ARRAY_SIZE]time.Time

var position int
var currentFloor int

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
		if orders[position] != 0 {
			break
		}
	}
	// move back to appearance of the floor
	for i := 0; i < ORDERS_ARRAY_SIZE; i++ {
		// avoid negative modulo
		if position == floor || position == ORDERS_ARRAY_SIZE-floor {
			return
		}
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

	// flush to disk // TODO

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
		if orders[index] != NONE {
			nextOrder = index
			break
		}
	}

	if nextOrder == -1 {
		return -1
	}

	// convert index to floor
	if nextOrder < ORDERS_ARRAY_SIZE/2 {
		return nextOrder

	}
	return ORDERS_ARRAY_SIZE - nextOrder
}

func GetDirection() int {

	next := GetNextOrder()

	fmt.Print("we are in floor ")
	fmt.Print(currentFloor)
	fmt.Print(" and want to go to ")
	fmt.Println(next)

	if next-currentFloor > 0 {
		return 1
	}
	if next-currentFloor < 0 {
		return -1
	}
	return 0
}
