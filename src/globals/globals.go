package globals

import (
	"network/udp"
)

var SignalChannel chan int = make(chan int, 5)

var MYID int = udp.GetMachineID()

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
