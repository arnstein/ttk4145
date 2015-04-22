package main

import (
	//"driver"
	//"fmt"
	//"iohandler"
	//	"network"
	//"queue"
	"errors"
	"fmt"
	"globals"
	"os"
	//"os/exec"
	"os/signal"
	"statemachine"
	//"watchdog"

	//"network/udp"
	//	"time"
	//	"net"
	//"iohandler"
)

func main() {

	//if len(os.Args) == 1 {

	//for {
	//watchdog.HbListener()
	//exec.Command("xterm", "-hold", "-e", "./src", "sdffsd").Start() // trollolo
	//}
	//}

	//go watchdog.UdpSend()
	//exec.Command("xterm", "-hold", "-e", "./src", "&").Start()

	defer statemachine.StateMachine()
	globals.SignalChannel <- globals.INIT
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			globals.CheckError(errors.New("YOMOMMA"))
			fmt.Println("LOLLOL", sig)
			//os.Exit(1)
		}
	}()
}
