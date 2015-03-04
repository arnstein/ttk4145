package udp

import (
	"fmt"
	"net"
	//	"strconv"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err.Error())
	}
}

func udpSend(sendChan <-chan []byte) {
	broadcastAddress, err := net.ResolveUDPAddr("udp", "129.241.187.255:20008")
	CheckError(err)
	sendSocket, err := net.DialUDP("udp", nil, broadcastAddress)
	CheckError(err)
	for {
		data := <-sendChan
		_, err := sendSocket.Write(data)
		CheckError(err)
		//	fmt.Println("sended " + strconv.Itoa(num) + " bytes with " + data)
	}
}

func udpRecv(recvChan chan<- []byte) {
	localAddress, err := net.ResolveUDPAddr("udp", ":20008")
	CheckError(err)
	recvSocket, err := net.ListenUDP("udp", localAddress)
	CheckError(err)
	var data []byte = make([]byte, 1500)
	for {
		length, addr, err := recvSocket.ReadFromUDP(data[0:])
		CheckError(err)
		fmt.Print("Message from  ")
		fmt.Print(addr.IP)
		fmt.Print(": ")
		recvChan <- data[:length]
	}
}

func UdpInit(sendChan <-chan []byte, recvChan chan<- []byte) {
	go udpSend(sendChan)
	go udpRecv(recvChan)
}
