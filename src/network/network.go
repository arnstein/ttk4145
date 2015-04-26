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

var activeOrderRequest [queue.ORDERS_ARRAY_SIZE]int
var costsOfOrders [queue.ORDERS_ARRAY_SIZE]responseData

var mess Message

var MessageDataChan = make(chan []int)
var sendChan = make(chan []byte)

type Message struct {
	MachineAddress int
	MessageType    int
	Data           []int
}

type responseData struct {
	Ip   int
	Cost int
}

func InitializeCostsOfOrders() {
	for i := 0; i < queue.ORDERS_ARRAY_SIZE; i++ {
		costsOfOrders[i].Ip = 260
		costsOfOrders[i].Cost = 42
	}
}

func NetworkInit() {
	receiveChan := make(chan []byte)
	udp.UdpInit(sendChan, receiveChan)
	go receiveMessage(receiveChan)
	go NewRequest()
}

func receiveMessage(receiveChan <-chan []byte) Message {
	for {
		receivedData := <-receiveChan
		decoded := decodeJSON(receivedData)
		parseMessage(decoded)
	}
}

func CheckBackupTimeouts() {
	for i := 0; i < queue.ORDERS_ARRAY_SIZE; i++ {
		queue.OrderBackup[i] = time.Unix(0, 0)
	}

	for {
		for i := 0; i < queue.ORDERS_ARRAY_SIZE; i++ {
			if queue.OrderBackup[i] == time.Unix(0, 0) {
				continue
			}
			if time.Since(queue.OrderBackup[i]) > 30*time.Second {
				floor, direction := queue.IndexToFloorAndDirection(i)
				globals.NewRequest <- [2]int{floor, direction}
				queue.OrderBackup[i] = time.Unix(0, 0)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func NewRequest() {
	var requestData [2]int
	for {
		requestData = <-globals.NewRequest
		floor := requestData[0]
		direction := requestData[1]
		message := Message{MachineAddress: globals.MYID, MessageType: ORDER, Data: []int{floor, direction}}
		sendChan <- encodeJSON(message)
	}
}

func RequestServed(floor int, direction int) {
	message := Message{MachineAddress: globals.MYID, MessageType: ORDERSERVED, Data: []int{floor, direction}}
	sendChan <- encodeJSON(message)
}

func putNewCost(cost int, ip int, index int) {
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
	orderIndex := queue.FloorAndDirToIndex(floor, direction)

	if activeOrderRequest[orderIndex] == 1 {
		return
	}
	activeOrderRequest[orderIndex] = 1

	cost := queue.CalculateCost(floor, direction)
	message := Message{MachineAddress: globals.MYID, MessageType: COSTORDER, Data: []int{floor, direction, cost}}
	sendChan <- encodeJSON(message)

	time.Sleep(1000 * time.Millisecond)
	if costsOfOrders[orderIndex].Ip == globals.MYID {
		queue.AddToQueue(floor, direction, queue.GLOBAL)
	} else {
		queue.AddToBackupQueue(floor, direction)
	}
	costsOfOrders[orderIndex].Cost = 42
	costsOfOrders[orderIndex].Ip = 260

	globals.LightsChannel <- [3]int{direction, floor, 1}
	activeOrderRequest[orderIndex] = 0
}

func encodeJSON(mess Message) []byte {
	byteArray, err := json.Marshal(mess)
	globals.CheckError(err)
	return byteArray
}

func decodeJSON(rawMessage []byte) Message {
	var message Message
	err := json.Unmarshal(rawMessage, &message)
	globals.CheckError(err)
	return message
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
		go handleNewRequest(message.Data[0], message.Data[1])
	case COSTORDER:
		fmt.Println("\t MessageType: CostOrder")
		index := queue.FloorAndDirToIndex(message.Data[0], message.Data[1])
		putNewCost(message.Data[2], message.MachineAddress, index)
		fmt.Print("\t Cost: ")
		fmt.Print(message.Data[2])
		fmt.Println()
	case ORDERSERVED:
		fmt.Println("\t MessageType: OrderServed")
		fmt.Print("\t Floor: ")
		fmt.Print(message.Data[0])
		fmt.Print(" Dir: ")
		fmt.Print(message.Data[1])
		fmt.Println()

		if message.MachineAddress == globals.MYID {
			queue.RemoveFromQueue(message.Data[0], message.Data[1])
		}
		queue.RemoveFromBackupQueue(message.Data[0], message.Data[1])
		globals.LightsChannel <- [3]int{message.Data[1], message.Data[0], 0}
	}
}
