package main

import (
	"errors"
	"globals"
	"os"
	"os/signal"
	"statemachine"
)

func main() {
	defer statemachine.StateMachine()
	globals.SignalChannel <- globals.INIT
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			globals.CheckError(errors.New("CTRL+C pressed"))
		}
	}()
}
