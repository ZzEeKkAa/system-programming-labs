package main

import (
	"fmt"
	"io"
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
	rand.Seed(time.Now().Unix())
	service := ":1234"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	defer listener.Close()

	for {
		func() {
			conn, err := listener.Accept()
			if err != nil {
				return
			}

			fmt.Println("Got connection")

			data, err := ReadAll(conn)
			checkError(err)

			fmt.Printf("Got data: %s\n", string(data))

			time.Sleep(time.Second * time.Duration(rand.Int31n(5)+5))

			fmt.Println("Returned data")

			conn.Write(data)
			conn.Close()
		}()
	}
}

func ReadAll(r io.Reader) ([]byte, error) {
	var buff = make([]byte, 10)

	n, err := r.Read(buff)
	if err != nil {
		return nil, err
	}

	buff = buff[:n]

	fmt.Println(n, string(buff))

	return buff, nil
}
