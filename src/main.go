package main

import (
	//"controlloop"
	"fmt"
	//	"math/rand"
	"network"
)

func main() {
	network.NetworkInit()
	arrived := make(chan int)
	//orders := make(chan int)
	//c := make(chan os.Signal, 1)
	//r := rand.New(rand.NewSource(99))

	//controlloop.InitCtrl(arrived, orders)
	fmt.Println("int done")
	for {
		<-arrived
	}
	//for {

	//		randnr := r.Intn(4)
	//		fmt.Println(randnr + 1)
	//		orders <- randnr
	//		<-arrived
	//	}
}
