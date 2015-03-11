package network

import (
	"encoding/json"
	"fmt"
	"network/udp"
	"time"
)

const (
	ORDER       = 0
	HEARTBEAT   = 1
	COSTORDER   = 2
	TAKEORDER   = 3
	ORDERSERVED = 4

	FLOOR     = 0
	DIRECTION = 1
	COST      = 2
)

type Message struct {
	messageType int
	data        string
}

	var message Message
	message.messageType = 200
	message.data = "lol"
func NetworkInit() {
	sendChan := make(chan []byte)
	receiveChan := make(chan []byte)
	udp.UdpInit(sendChan, receiveChan)
	go sendMessage(sendChan)
	go receiveMessage(receiveChan)
}

func receiveMessage(receiveChan <-chan []byte) Message {
	//heartbeatTime := time.Now()
	for {
		receivedData := <-receiveChan
		//heartbeatTime = time.Now()
		decoded := decodeJSON(receivedData)
		fmt.Println(decoded.data)
	}

}

func sendMessage(sendChan chan<- []byte) {
	for {
		time.Sleep(1 * time.Second)
		sendChan <- encodeJSON(message)
	}
}

func encodeJSON(mess Message) []byte {
	b, err := json.Marshal(mess)
	udp.CheckError(err)
	return b
}

func decodeJSON(mess []byte) Message {
	var me Message

	err := json.Unmarshal(mess, &me)
	udp.CheckError(err)
	return me
}
