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

// dummy comment

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

func InitializeCostsOfOrders() {
	for i := 0; i < queue.ORDERS_ARRAY_SIZE; i++ {
		costsOfOrders[i].Ip = 260
		costsOfOrders[i].Cost = 42
	}
}

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

func CheckBackupTimeouts() {
	for i := 0; i < queue.ORDERS_ARRAY_SIZE; i++ {
		queue.OrderBackup[i] = time.Unix(0, 0)
	}

	for {
		for i := 0; i < queue.ORDERS_ARRAY_SIZE; i++ {
			if queue.OrderBackup[i] == time.Unix(0, 0) {
				continue
			}
			if time.Since(queue.OrderBackup[i]) > 1*time.Minute {
				fmt.Println("now I shoud resend something")
				floor, direction := queue.IndexToFloorAndDirection(i)
				NewRequest(floor, direction)
				queue.OrderBackup[i] = time.Unix(0, 0)
			}
		}
		time.Sleep(1 * time.Second)
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
	orderIndex := queue.FloorAndDirToIndex(floor, direction)
	if activeOrderRequest[orderIndex] == 1 {
		return
	}
	activeOrderRequest[orderIndex] = 1
	cost := queue.CalculateCost(floor, direction)
	message := Message{MachineAddress: globals.MYID, MessageType: COSTORDER, Data: []int{floor, direction, cost}}
	sendChan <- encodeJSON(message)

	//fmt.Println("IT GOES HERE RIGHT!!!!!!!")
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("inside handlenewreq")
	fmt.Println(costsOfOrders)
	if costsOfOrders[orderIndex].Ip == globals.MYID {
		fmt.Println("inside this if")
		queue.AddToQueue(floor, direction, queue.GLOBAL)
	} else {
		//fmt.Println("inside this ELIF LOL")
		queue.AddToBackupQueue(floor, direction)
	}
	//fmt.Println("DOESN'T IT GO HERE???")
	//time.Sleep(1000 * time.Millisecond)
	costsOfOrders[orderIndex].Cost = 42
	costsOfOrders[orderIndex].Ip = 260 // do we need this?
	driver.SetOutsideLamp(floor, direction)
	activeOrderRequest[orderIndex] = 0
}

func encodeJSON(mess Message) []byte {
	b, err := json.Marshal(mess)
	globals.CheckError(err)
	return b
}

func decodeJSON(mess []byte) Message {
	var me Message
	err := json.Unmarshal(mess, &me)
	globals.CheckError(err)
	return me
}

func parseMessage(message Message) {
	//queue.PrintQueue()
	fmt.Println(costsOfOrders)
	fmt.Println(activeOrderRequest)
	queue.PrintQueue()
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
	case HEARTBEAT:
		fmt.Println("\t MessageType: Heartbeat")
	case COSTORDER:
		fmt.Println("\t MessageType: Costorder")
		index := queue.FloorAndDirToIndex(message.Data[0], message.Data[1])
		putNewCost(message.Data[2], message.MachineAddress, index)
		fmt.Print(message.Data[2])
		fmt.Print(" was the cost")
		fmt.Println()
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
