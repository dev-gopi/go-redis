package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {

	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to server")

	resp := "*1\r\n$4\r\nPING\r\n"

	_, err = conn.Write([]byte(resp))
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(conn)

	reply, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	fmt.Println("Server Response:", reply)
}
