package udp

import (
	"globals"
	"net"
	"strconv"
	"strings"
    "time"
)

var sendSocketEnabled := 0
var receiveSocketEnabled := 0

func retryInit(sendChan <-chan []byte, receiveChan chan<- []byte) {
    time.Sleep(1*time.Minute)
    UdpInit(sendChan, receiveChan)
}

func UdpInit(sendChan <-chan []byte, receiveChan chan<- []byte) {
	globals.MYID = getMachineID()
	broadcastAddress, err := net.ResolveUDPAddr("udp", "129.241.187.255:20008")
    if !sendSocketEnabled {
        sendSocket, err := net.DialUDP("udp", nil, broadcastAddress)
        if err == nil {
            sendSocketEnabled = 1
        }
	    go udpSend(sendChan)
    }

	localAddress, err := net.ResolveUDPAddr("udp", ":20008")
    if !reciveSocketEnabled {
	    receiveSocket, err := net.ListenUDP("udp", localAddress)
        if err == nil {
            receiveSocketEnabled = 1
	        go udpReceive(receiveChan)
        }
    }
    if !(sendSocketEnabled && receiveSocketEnabled){
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

func udpSend(sendChan <-chan []byte) {
	for {
		data := <-sendChan
        if sendSocketEnabled {
		    _, err := sendSocket.Write(data)
		    globals.CheckError(err)
        }
	}
}

func udpReceive(receiveChan chan<- []byte,) {
	for {
	    var data []byte = make([]byte, 1500)
	    length, _, err := receiveSocket.ReadFromUDP(data[0:])
	    globals.CheckError(err)
	    receiveChan <- data[:length]
	}
}
