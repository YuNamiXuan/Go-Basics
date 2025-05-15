package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	serverAddr, err := net.ResolveUDPAddr("udp", ":8081")
	if err != nil {
		fmt.Println("Resolve address failed: ", err)
		return
	}

	localAddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		fmt.Println("Resolve local address failed: ", err)
		return
	}

	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		fmt.Println("Connect to server failed: ", err)
		return
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Input failed: ", err)
			continue
		}

		if strings.TrimSpace(message) == "exit" {
			fmt.Println("Client exit.")
			return
		}

		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Send message failed: ", err)
			return
		}

		buffer := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Read response from server failed: ", err)
			return
		}
		fmt.Println("Server response: ", string(buffer[:n]))
	}
}
