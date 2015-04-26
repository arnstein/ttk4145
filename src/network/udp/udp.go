package udp

import (
	"fmt"
	"globals"
	"net"
	"strconv"
	"strings"
	"time"
)

var sendSocketEnabled int = 0
var receiveSocketEnabled int = 0

func retryInit(sendChan <-chan []byte, receiveChan chan<- []byte) {
	fmt.Println("Network not work. Retrying to connect in 10 seconds")
	time.Sleep(10 * time.Second)
	UdpInit(sendChan, receiveChan)
}

func UdpInit(sendChan <-chan []byte, receiveChan chan<- []byte) {
	globals.MYID = getMachineID()
	broadcastAddress, _ := net.ResolveUDPAddr("udp", "129.241.187.255:20008")
	if sendSocketEnabled == 0 {
		sendSocket, err := net.DialUDP("udp", nil, broadcastAddress)
		if err == nil {
			sendSocketEnabled = 1
		}
		go udpSend(sendChan, sendSocket)
	}

	localAddress, _ := net.ResolveUDPAddr("udp", ":20008")
	if receiveSocketEnabled == 0 {
		receiveSocket, err := net.ListenUDP("udp", localAddress)
		if err == nil {
			receiveSocketEnabled = 1
			go udpReceive(receiveChan, receiveSocket)
		}
	}
	if sendSocketEnabled+receiveSocketEnabled != 2 {
		go retryInit(sendChan, receiveChan)
	}
}

func getMachineID() int {
	var localAddr string
	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		if addrs, err := inter.Addrs(); err == nil {
			for _, addr := range addrs {
				if inter.Name == "eth0" && strings.Contains(addr.String(), "129") {
					localAddr = addr.String()
				}
			}
		}
	}
	localAddr = strings.Split(localAddr, "/")[0]
	localAddr = strings.Split(localAddr, ".")[3]
	localID, _ := strconv.Atoi(localAddr)
	return localID
}

func udpSend(sendChan <-chan []byte, sendSocket *net.UDPConn) {
	if sendSocketEnabled == 0 {
		for sendSocketEnabled == 0 {
			<-sendChan
		}
		return
	}
	for {
		data := <-sendChan
		_, err := sendSocket.Write(data)
		globals.CheckError(err)
	}
}

func udpReceive(receiveChan chan<- []byte, receiveSocket *net.UDPConn) {
	for {
		var data []byte = make([]byte, 1500)
		length, _, err := receiveSocket.ReadFromUDP(data[0:])
		globals.CheckError(err)
		receiveChan <- data[:length]
	}
}
