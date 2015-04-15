package network

import (
	"driver"
	"encoding/json"
	"fmt"
	"globals"
	"network/udp"
	"queue"
	"sync"
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

var activeOrderRequest [queue.ORDERS_ARRAY_SIZE]int // find out if globals or queue
var mutex = &sync.Mutex{}

type Message struct {
	MachineAddress int
	MessageType    int
	Data           []int
}

type responseData struct {
	Ip   int
	Cost int
}

var costsOfOrders [queue.ORDERS_ARRAY_SIZE]responseData

var mess Message
var MessageDataChan = make(chan []int)
var sendChan = make(chan []byte)

func NetworkInit() {
	receiveChan := make(chan []byte)
	udp.UdpInit(sendChan, receiveChan)
	go receiveMessage(receiveChan)
}

func receiveMessage(receiveChan <-chan []byte) Message {
	//heartbeatTime := time.Now()
	for {
		receivedData := <-receiveChan
		//heartbeatTime = time.Now()
		//fmt.Print("received data: ")
		//fmt.Println(receivedData)
		decoded := decodeJSON(receivedData)
		//fmt.Print("decoded data:  ")
		//fmt.Println(decoded)

		parseMessage(decoded)
	}

}

func NewRequest(floor int, direction int) {
	message := Message{MachineAddress: globals.MYID, MessageType: ORDER, Data: []int{floor, direction}}
	sendChan <- encodeJSON(message)
}

func RequestServed(floor int, direction int) {
	message := Message{MachineAddress: globals.MYID, MessageType: ORDERSERVED, Data: []int{floor, direction}}
	sendChan <- encodeJSON(message)
}

func putNewCost(cost int, ip int, index int) {
	//need syncronisatioin!
	if costsOfOrders[index].Cost < cost {
		return
	}
	if costsOfOrders[index].Cost == cost &&
		costsOfOrders[index].Ip < ip {
		return
	}

	costsOfOrders[index].Ip = ip
	costsOfOrders[index].Cost = cost

}

func handleNewRequest(floor int, direction int) {
	lowest := globals.MYID
	orderIndex := queue.FloorAndDirToIndex(floor, direction)
	if activeOrderRequest[orderIndex] == 1 {
		return
	}
	activeOrderRequest[orderIndex] = 1

	cost := queue.CalculateCost(floor, direction)
	message := Message{MachineAddress: globals.MYID, MessageType: COSTORDER, Data: []int{floor, direction, cost}}
	sendChan <- encodeJSON(message)

	time.Sleep(1000 * time.Millisecond)
	// sjekk channel
	if lowest == globals.MYID {
		queue.AddToQueue(floor, direction, queue.GLOBAL)
	} else {
		queue.AddToBackupQueue(floor, direction)
	}
	driver.SetOutsideLamp(floor, direction)
	//set lights
	activeOrderRequest[orderIndex] = 0
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
	fmt.Print("New message from MachineID: ")
	fmt.Print(message.MachineAddress)
	fmt.Println()
	switch message.MessageType {
	case ORDER:
		fmt.Println("\t MessageType: Order")
		fmt.Print("\t Floor: ")
		fmt.Print(message.Data[0])
		fmt.Print(" Dir: ")
		fmt.Print(message.Data[1])
		fmt.Println()
		handleNewRequest(message.Data[0], message.Data[1])
	case HEARTBEAT:
		fmt.Println("\t MessageType: Heartbeat")
	case COSTORDER:
		fmt.Println("\t MessageType: Costorder")
		// add to channel
	case TAKEORDER:
		fmt.Println("\t MessageType: Take order")
	case ORDERSERVED:
		fmt.Println("\t MessageType: Order served")
		fmt.Print("\t Floor: ")
		fmt.Print(message.Data[0])
		fmt.Print(" Dir: ")
		fmt.Print(message.Data[1])
		fmt.Println()
		if message.MachineAddress == globals.MYID {
			queue.RemoveFromQueue(message.Data[0], message.Data[1])
		}
		queue.RemoveFromBackupQueue(message.Data[0], message.Data[1])
		driver.ClearOutsideLamp(message.Data[0], message.Data[1])

	}
}
