package network

import (
	"encoding/json"
	"fmt"
	"globals"
	"network/udp"
	"queue"
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
var messageChan = make(chan Message)

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

func NewRequest(floor int, direction int) {
	message := Message{MachineAddress: globals.MYID, MessageType: ORDER, Data: []int{floor, direction}}
	messageChan <- message
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
	lowest := 154
	orderIndex := queue.FloorAndDirToIndex(floor, direction)
	if activeOrderRequest[orderIndex] == 1 {
		return
	}
	activeOrderRequest[orderIndex] = 1
	// send egen cost
	// make new message with globals.MYID, calculatecost, floor and direction
	message := Message{MachineAddress: globals.MYID, MessageType: COSTORDER, Data: []int{floor, direction}}
	queue.CalculateCost(floor, direction)

	time.Sleep(1000 * time.Millisecond)
	// sjekk channel
	if lowest == globals.MYID {
		queue.AddToQueue(floor, direction, queue.GLOBAL)

	}
}
func sendMessage(sendChan chan<- []byte) {
	//	mess := Message{MachineAddress: udp.GetMachineID(), MessageType: ORDERSERVED, Data: []int{1, 2, 3}}
	var message Message
	for {
		message = <-messageChan
		//	time.Sleep(1 * time.Second)
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
		// go handleNowReq
	case HEARTBEAT:
		fmt.Println("\t MessageType: Heartbeat")
	case COSTORDER:
		fmt.Println("\t MessageType: Costorder")
		// add to channel
	case TAKEORDER:
		fmt.Println("\t MessageType: Take order")
	case ORDERSERVED:
		fmt.Println("\t MessageType: Order served")

	}
}
