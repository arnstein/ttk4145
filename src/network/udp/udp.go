package udp

import (
	"fmt"
	"net"
	//	"strconv"
	//"strings"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("error")
		fmt.Println(err.Error())
	}
}

/*
func GetLocalAddr() {
	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		if addrs, err := inter.Addrs(); err == nil {
			for _, addr := range addrs {
				if inter.Name == "eth0" && strings.Contains(addr.String(), "129") {
					fmt.Println(addr)
				}
			}
		}
	}

}

*/
func udpSend(sendChan <-chan []byte) {
	broadcastAddress, err := net.ResolveUDPAddr("udp", "129.241.187.255:20008")
	CheckError(err)
	sendSocket, err := net.DialUDP("udp", nil, broadcastAddress)
	CheckError(err)
	for {
		data := <-sendChan
		_, err := sendSocket.Write(data)
		CheckError(err)
	}
}

func udpReceive(receiveChan chan<- []byte) {
	localAddress, err := net.ResolveUDPAddr("udp", ":20008")
	CheckError(err)
	receiveSocket, err := net.ListenUDP("udp", localAddress)
	CheckError(err)
	var data []byte = make([]byte, 1500)
	for {
		length, addr, err := receiveSocket.ReadFromUDP(data[0:])
		CheckError(err)
		fmt.Print("Message from  ")
		fmt.Print(addr.IP)
		fmt.Print(": ")
		receiveChan <- data[:length]
	}
}

func UdpInit(sendChan <-chan []byte, receiveChan chan<- []byte) {
	go udpSend(sendChan)
	go udpReceive(receiveChan)
}
