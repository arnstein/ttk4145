package main

import (
	"net"
)

func main() {
	client, error  := net.Dial("tcp","129.241.187.136:33546"); 
	if error != nil{
		println("dial");
		println(error.Error());
		return;
	}

	sendData := []byte("Connect to: 129.241.187.161:33333\000");

	println("setup complete");

	reply := make([]byte, 1024);
	_, error = client.Read(reply);
	if error != nil{
		println(error.Error());
		return;
	}
	println("reply from server=", string(reply));

	_, error = client.Write(sendData);
	if error != nil{
		println(error.Error());
		return;
	}
	println("send  complete");
}
