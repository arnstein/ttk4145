package udp

import (
	"fmt"
	"net"
	//	"strconv"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func udpSend(sendChan <-chan string) {
	broadcastAddress, err := net.ResolveUDPAddr("udp", "129.241.187.255:20008")
	checkError(err)
	sendSocket, err := net.DialUDP("udp", nil, broadcastAddress)
	checkError(err)
	for {
		data := <-sendChan
		_, err := sendSocket.Write([]byte(data))
		checkError(err)
		//	fmt.Println("sended " + strconv.Itoa(num) + " bytes with " + data)
	}
}

func udpRecv(recvChan chan<- string) {
	localAddress, err := net.ResolveUDPAddr("udp", ":20008")
	checkError(err)
	recvSocket, err := net.ListenUDP("udp", localAddress)
	checkError(err)
	var data []byte = make([]byte, 1500)
	for {
		length, addr, err := recvSocket.ReadFromUDP(data[0:])
		checkError(err)
		fmt.Print("Message from  ")
		fmt.Print(addr.IP)
		fmt.Print(": ")
		recvChan <- string(data[:length])
	}
}

func UdpInit(sendChan <-chan string, recvChan chan<- string) {
	go udpSend(sendChan)
	go udpRecv(recvChan)
}
