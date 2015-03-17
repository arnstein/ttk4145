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
	MessageType int
	Data        []int
}

var mess Message

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
		parseMessage(decoded)
	}

}

func sendMessage(sendChan chan<- []byte) {
	mess := Message{MessageType: ORDERSERVED, Data: []int{1, 2, 3}}
	for {
		for i := 0; i < 5; i++ {
			mess.MessageType = i
			time.Sleep(1 * time.Second)
			sendChan <- encodeJSON(mess)
			fmt.Println("Sent message")
		}
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

func parseMessage(message Message) {
	switch message.MessageType {
	case ORDER:
		fmt.Println("MessageType: Order")
	case HEARTBEAT:
		fmt.Println("MessageType: Heartbeat")
	case COSTORDER:
		fmt.Println("MessageType: Costorder")
	case TAKEORDER:
		fmt.Println("MessageType: Take order")
	case ORDERSERVED:
		fmt.Println("MessageType: Order served")

	}
}
