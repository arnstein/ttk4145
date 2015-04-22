package watchdog

import (
	"fmt"
	"net"
	"time"

	"globals"
)

func UdpSend() {
	sendAddress, err := net.ResolveUDPAddr("udp", "localhost:24200")
	localAddress, err := net.ResolveUDPAddr("udp", "localhost:24201")
	globals.CheckError(err)
	sendSocket, err := net.ListenUDP("udp", localAddress)
	globals.CheckError(err)
	for {
		time.Sleep(500 * time.Millisecond)
		_, err := sendSocket.WriteToUDP([]byte("Sent"), sendAddress)
		globals.CheckError(err)
	}
}

func udpRecv() bool {

	localAddress, err := net.ResolveUDPAddr("udp", "127.0.0.1:24200")
	globals.CheckError(err)

	recvSocket, err := net.ListenUDP("udp", localAddress)
	globals.CheckError(err)
	data := make([]byte, 1024)
	recvSocket.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	length, _, err := recvSocket.ReadFromUDP(data)
	if err != nil {
		return false
	}
	globals.CheckError(err)

	fmt.Println(length)

	return true
}

func HbListener() {

	// sla pa listener socket
	localAddress, err := net.ResolveUDPAddr("udp", "127.0.0.1:24200")
	globals.CheckError(err)

	recvSocket, err := net.ListenUDP("udp", localAddress)
	globals.CheckError(err)
	data := make([]byte, 1024)

	lastHb := time.Now()

	for time.Since(lastHb) <= 2*time.Second {
		recvSocket.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		_, _, err := recvSocket.ReadFromUDP(data)
		if err == nil {
			fmt.Println("recievd ")
			lastHb = time.Now()
		}
	}
	// sla av listener socket
	recvSocket.Close()
}
