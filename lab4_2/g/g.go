package main

import (
	"math/rand"
	"net"
	"time"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	service := ":4321"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		time.Sleep(time.Second * time.Duration(rand.Int31n(5)+5))

		conn.Write([]byte("false"))
		conn.Close()
	}
}
