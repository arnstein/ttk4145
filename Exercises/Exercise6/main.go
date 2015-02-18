package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func checkError(err error) {

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

}

func udpSend() {
	sendAddress, err := net.ResolveUDPAddr("udp", "localhost:20008")
	localAddress, err := net.ResolveUDPAddr("udp", "localhost:20009")
	checkError(err)
	sendSocket, err := net.ListenUDP("udp", localAddress)
	checkError(err)
	for {
		time.Sleep(500 * time.Millisecond)
		_, err := sendSocket.WriteToUDP([]byte("Sent"), sendAddress)
		checkError(err)
	}
}

func udpRecv() bool {

	localAddress, err := net.ResolveUDPAddr("udp", "127.0.0.1:20008")
	checkError(err)

	recvSocket, err := net.ListenUDP("udp", localAddress)
	checkError(err)
	data := make([]byte, 1024)
	recvSocket.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	length, _, err := recvSocket.ReadFromUDP(data)
	if err != nil {
		return false
	}
	checkError(err)

	fmt.Println(length)

	return true
}

func hbListener() {

	// sla pa listener socket
	localAddress, err := net.ResolveUDPAddr("udp", "127.0.0.1:20008")
	checkError(err)

	recvSocket, err := net.ListenUDP("udp", localAddress)
	checkError(err)
	data := make([]byte, 1024)
	recvSocket.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

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

func openFile() *os.File {
	f, err := os.OpenFile("dat1", os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("create new file")
		f, err = os.Create("dat1")
		checkError(err)
		_, err = f.Write([]byte("0\n"))
		checkError(err)
	} else {
		fmt.Println("read file")
	}

	return f

}

func readLastNumber() int {
	buff, err := exec.Command("tail", "-1", "dat1").Output()
	checkError(err)
	s := string(buff)
	s = strings.TrimSuffix(s, "\n")
	number, err := strconv.Atoi(s)
	checkError(err)
	return number

}
func main() {

	//backup
	hbListener()

	//new primery
	fmt.Println("We are new Primary")

	// heartbeatsender
	go udpSend()

	err := exec.Command("xterm", "go", "run", "main.go").Start()
	checkError(err)

	f := openFile()
	var number int = readLastNumber()

	//primary
	for {
		d1 := []byte(strconv.Itoa(number) + "\n")
		_, err := f.Write(d1)
		checkError(err)
		fmt.Println(number)
		number++

		time.Sleep(1 * time.Second)

	}

}
