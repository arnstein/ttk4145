package queue

import (
	"fmt"
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
var position int

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

}

func removeFromQueue() {

	orders[position] = NONE

}

func CalculateCost(floor int, dir int) int {

	endIndex := FloorAndDirToIndex(floor, dir)

	cost := 0

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
