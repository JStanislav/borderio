package main

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	fmt.Printf("New connection from %s\n", conn.RemoteAddr().String())
	input := make([]byte, 1024)
	n, err := conn.Read(input)
	if err != nil {
		fmt.Printf("[ERROR] reading: %v\n", err)
	}
	fmt.Printf("Reading n: %d\n", n)
	fmt.Printf("Message: %s\n", input[:n])

	conn.Write([]byte("Hello from the server"))

	conn.Close()
}

func main() {
	fmt.Println("Hello World")

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}
}
