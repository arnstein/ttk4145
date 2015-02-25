package main

import (
	"controlloop"
	"fmt"
	"math/rand"
	"network"
	"os"
	"os/signal"
)

func main() {
	network.NetworkInit()
	arrived := make(chan int)
	orders := make(chan int)
	c := make(chan os.Signal, 1)
	r := rand.New(rand.NewSource(99))

	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Println(sig)
			os.Exit(1)
		}
	}()
	controlloop.InitCtrl(arrived, orders)
	fmt.Println("int done")

	for {

		randnr := r.Intn(4)
		fmt.Println(randnr + 1)
		orders <- randnr
		<-arrived
	}
}
