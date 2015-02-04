package main

import (
	"fmt"
	"time"
	"udp"
)

func main() {

	sendChan := make(chan string)
	recvChan := make(chan string)

	udp.UdpInit(sendChan, recvChan)
	fmt.Println("init done")

	for {
		sendChan <- "Hajime!"
		time.Sleep(1000 * time.Millisecond)
	}

}
