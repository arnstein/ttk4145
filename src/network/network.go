package network

import (
	"encoding/json"
	"fmt"
	"network/udp"
	"time"
)

type Message struct {
	Identifier  string
	MessageType string
	Time        int64
}

var m Message

func NetworkInit() {
	sendChan := make(chan []byte)
	receiveChan := make(chan []byte)
	udp.UdpInit(sendChan, receiveChan)
	go sendHeartBeat(sendChan)
	go receiver(receiveChan)
}

func receiver(receiveChan <-chan []byte) Message {
	//heartbeatTime := time.Now()
	for {
		receivedData := <-receiveChan
		//heartbeatTime = time.Now()
		decoded := decodeJSON(receivedData)
		fmt.Println(decoded)
	}

}

func UpdateHeartBeatInfo(name string, messageType string) {
	m.Identifier = name
	m.MessageType = messageType
	m.Time = time.Now().Unix()
}

func sendHeartBeat(sendChan chan<- []byte) {
	for {
		UpdateHeartBeatInfo("Bob", "heartbeat")
		time.Sleep(1 * time.Second)
		sendChan <- encodeJSON(m)
	}
}

func encodeJSON(mess Message) []byte {
	b, err := json.Marshal(mess)
	udp.CheckError(err)
	return b
}

func decodeJSON(mess []byte) Message {
	var me Message

	eror := json.Unmarshal(mess, &me)
	udp.CheckError(eror)

	return me
}
