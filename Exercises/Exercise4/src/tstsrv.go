package main

import (
	"udp"
)

func tstsrv() {

	udpInit()
	buf := make([]byte, 1024)

	udpRecv(buf)

}
