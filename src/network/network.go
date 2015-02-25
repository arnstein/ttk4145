package network

import (
	"fmt"
	"network/udp"
	"time"
)

func NetworkInit() {
	sendChan := make(chan string)
	recvChan := make(chan string)
	udp.UdpInit(sendChan, recvChan)
	go sendHeartBeat(sendChan)
	go receiver(recvChan)
}

func receiver(recvChan <-chan string) {
	//heartbeatTime := time.Now()
	for {
		receivedData := <-recvChan
		//heartbeatTime = time.Now()
		fmt.Println("Received " + receivedData)
	}

}

func UpdateHeartBeatInfo(newHB string) {

}

func sendHeartBeat(sendChan chan<- string) {
	for {
		time.Sleep(1 * time.Second)
		sendChan <- "blubb"
	}
}
