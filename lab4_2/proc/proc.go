package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

var (
	fServer = flag.String("f_serv", ":1234", "f server address")
	gServer = flag.String("g_serv", ":4321", "g server address")
)

func main() {
	flag.Parse()

	fTcpAddr, err := net.ResolveTCPAddr("tcp4", *fServer)
	checkError(err)
	fConn, err := net.DialTCP("tcp", nil, fTcpAddr)
	checkError(err)

	_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	checkError(err)
	result, err := ioutil.ReadAll(conn)
	checkError(err)
	fmt.Println(string(result))
	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
