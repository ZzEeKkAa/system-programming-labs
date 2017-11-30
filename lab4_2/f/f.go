package main

import (
	"io/ioutil"
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
	service := ":1234"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		func() {
			conn, err := listener.Accept()
			if err != nil {
				return
			}

			data, err := ioutil.ReadAll(conn)
			checkError(err)

			time.Sleep(time.Second * time.Duration(rand.Int31n(5)+5))

			conn.Write([]byte("true"))
			conn.Close()
		}()
	}
}
