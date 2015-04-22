package globals

import (
	"fmt"
)

var SignalChannel chan int = make(chan int, 5)

var MYID int

const (
	NUM_FLOORS = 4
)
const (
	INIT         = 0
	MOVEORDER    = 1
	TIMEROUT     = 2
	EMPTYQUEUE   = 3
	FLOORREACHED = 4
	CHECKORDER   = 5
)

const (
	STOP = 0
	UP   = 1
	DOWN = -1
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("error")
		fmt.Println(err.Error())
	}
}
