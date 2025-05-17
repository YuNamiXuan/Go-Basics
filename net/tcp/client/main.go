package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Create a socket and connect to server
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println("Conntect to server failed.")
		return
	}
	defer conn.Close()

	// Create a buffered reader for standard input (user input)
	reader := bufio.NewReader(os.Stdin)

	// Main client loop
	for {
		fmt.Print(">")

		// Read user input
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Input message failed.")
			return
		}

		// Send data to server
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Send message failed.")
		}

		if strings.TrimSpace(message) == "exit" {
			conn.Write([]byte("exit\n"))
			fmt.Printf("Client %s exit.\n", conn.LocalAddr().String())
			return
		}

		// Read server response
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Read response from server failed.")
			return
		}
		fmt.Print("Server response: " + response)
	}
}
