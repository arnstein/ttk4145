package udp

import (
	"fmt"
	"net"
	"strconv"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func udpSend(sendChan <-chan string) {

	broadcastAddress, err := net.ResolveUDPAddr("udp", "129.241.187.255:20008")
	checkError(err)
	localAddress, err := net.ResolveUDPAddr("udp", "129.241.187.161:20008")
	//sendSocket, err := net.DialUDP("udp", nil, broadcastAddress)
	tempSocket, err := net.DialUDP("udp4", nil, broadcastAddress)
	checkError(err)
	defer tempSocket.Close()
	sendSocket, err := net.ListenUDP("udp", localAddress)
	//	checkError(err)

	for {

		data := <-sendChan
		num, err := sendSocket.WriteToUDP([]byte(data), broadcastAddress)
		//num, err := sendSocket.WriteToUDP([]byte(data), &net.UDPAddr{IP: net.IP{129, 241, 187, 255}, Port: 20008})
		checkError(err)
		fmt.Println("sended " + strconv.Itoa(num) + " bytes with " + data)
	}
}

func udpRecv(Jackie chan<- string) {

	localAddress, err := net.ResolveUDPAddr("udp", "127.0.0.1:20008")
	checkError(err)

	recvSocket, err := net.ListenUDP("udp", localAddress)
	checkError(err)

	data := make([]byte, 7)

	for {

		length, _, err := recvSocket.ReadFromUDP(data)
		checkError(err)
		fmt.Println(string(data[:length]))
		//Jackie <- string(data[:length])
	}
}

func UdpInit(sendChan <-chan string, recvChan chan<- string) {

	go udpSend(sendChan)
	go udpRecv(recvChan)
}
