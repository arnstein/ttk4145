package main

import (
	"net"
)

func main() {
	server, error  := net.Listen("tcp",":33333"); 
	if error != nil{
		println(error.Error());
		return;
	}

	for{

		conn, error := server.Accept();
		if error != nil{
			println(error.Error());
			return;
		}

		go func(c net.Conn){
			sendData := []byte("bla\000");
			reply := make([]byte, 1024);
			_, error = c.Read(reply);
			_, error = c.Write(sendData);
			_, error = c.Read(reply);
			println("reply from server=", string(reply));
			c.Close();

		}(conn)
	}
	server.Close();
}
