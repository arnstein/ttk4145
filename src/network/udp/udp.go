package udp

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("error")
		fmt.Println(err.Error())
	}
}

func GetMachineID() int {
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

func udpSend(sendChan <-chan []byte) {
	broadcastAddress, err := net.ResolveUDPAddr("udp", "129.241.187.255:20008")
	CheckError(err)
	sendSocket, err := net.DialUDP("udp", nil, broadcastAddress)
	CheckError(err)
	for {
		data := <-sendChan
		//fmt.Print("Sending:   ")
		//fmt.Println(data)
		_, err := sendSocket.Write(data)
		CheckError(err)
	}
}

func udpReceive(receiveChan chan<- []byte) {
	localAddress, err := net.ResolveUDPAddr("udp", ":20008")
	CheckError(err)
	receiveSocket, err := net.ListenUDP("udp", localAddress)
	CheckError(err)
	for {
		var data []byte = make([]byte, 1500)
		length, _, err := receiveSocket.ReadFromUDP(data[0:])
		CheckError(err)
		//fmt.Print("Message from  ")
		//fmt.Print(addr.IP)
		//fmt.Print(": ")
		//tmpData := data[:length]
		receiveChan <- data[:length]
	}
}

func UdpInit(sendChan <-chan []byte, receiveChan chan<- []byte) {
	go udpSend(sendChan)
	go udpReceive(receiveChan)
}
