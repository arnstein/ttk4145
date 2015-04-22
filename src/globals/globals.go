package globals

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

var SignalChannel chan int = make(chan int, 5)
var NewRequest chan [2]int = make(chan [2]int)
var LightsChannel chan [3]int = make(chan [3]int)

var MYID int

const (
	NUM_FLOORS = 4

	INIT         = 0
	MOVEORDER    = 1
	TIMEROUT     = 2
	EMPTYQUEUE   = 3
	FLOORREACHED = 4
	CHECKORDER   = 5

	STOP = 0
	UP   = 1
	DOWN = -1

	BUTTON_COMMAND = 2
)

func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())

		cmd := exec.Command("xterm", "-hold", "-e", "./src")
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
		cmd.Start()

		os.Exit(1)
	}
}
