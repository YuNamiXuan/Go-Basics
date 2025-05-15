package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println("Conntect to server failed.")
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Input message failed.")
			return
		}

		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Send message successfully.")
		}

		if strings.TrimSpace(message) == "exit" {
			conn.Write([]byte("exit\n"))
			fmt.Printf("Client %s exit.\n", conn.LocalAddr().String())
			return
		}

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Read response from server failed.")
			return
		}
		fmt.Print("Server response: " + response)
	}
}
