package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func ConnectionHandler(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Client %s connect successfully.\n", conn.RemoteAddr().String())

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Receive message failed: ", err)
		}

		message = strings.TrimSpace(message)
		if message == "exit" {
			fmt.Printf("Client %s exit.\n", conn.RemoteAddr().String())
			return
		} else {
			fmt.Printf("Client %s send: %s\n", conn.RemoteAddr().String(), message)
		}

		response := fmt.Sprintln("Server receive successfully.")
		conn.Write([]byte(response))
	}
}

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Listen failed: ", err)
		return
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Accept failed: ", err)
			continue
		}
		go ConnectionHandler(conn)
	}
}
