package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
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

	defer fConn.Close()

	gTcpAddr, err := net.ResolveTCPAddr("tcp4", *gServer)
	checkError(err)
	gConn, err := net.DialTCP("tcp", nil, gTcpAddr)
	checkError(err)
	defer gConn.Close()

	fmt.Println("Put two boolean in string form (example: true false):")

	var a, b bool
	fmt.Scanf("%t %t\n", &a, &b)

	fmt.Printf("Your input: %t %t\n", a, b)

	fmt.Println("Writing f")
	fConn.Write([]byte(fmt.Sprintf("%t", a)))

	fmt.Println("Writing g")
	gConn.Write([]byte(fmt.Sprintf("%t", b)))

	fmt.Println("Getting answer f")
	fRes, err := ioutil.ReadAll(fConn)
	checkError(err)

	fmt.Println("Getting answer g")
	gRes, err := ioutil.ReadAll(gConn)
	checkError(err)

	aa, err := strconv.ParseBool(string(fRes))
	checkError(err)
	bb, err := strconv.ParseBool(string(gRes))
	checkError(err)

	fmt.Printf("f(%t) && g(%t)=%t\n", a, b, aa && bb)

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
