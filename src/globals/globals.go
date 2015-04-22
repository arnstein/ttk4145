package globals

import (
	"fmt"
	"os"
	"syscall"

	"os/exec"
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
		//exec.Command("./src &").Start() // trollolo
		//exec.command("xterm", "-hold", "-e", "go", "run", "main.go").start()
		cmd := exec.Command("xterm", "-hold", "-e", "./src")
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
		cmd.Start()
		//syscall.ForkExec("xterm", []string{"-hold", "-e", "./src"}, nil)
		os.Exit(1)
	}
}
