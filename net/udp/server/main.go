package main

import (
	"fmt"
	"net"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":8081")
	if err != nil {
		fmt.Println("Resolve address failed: ", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Listen failed: ", err)
		return
	}
	defer conn.Close()

	fmt.Println("UDP Server Running, Port: 8081...")

	buffer := make([]byte, 1024)
	for {
		n, clietnAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Read failed: ", err)
			continue
		}

		message := string(buffer[:n])
		fmt.Printf("Client %s send message: %s", clietnAddr.String(), message)

		response := []byte("Receive message successfully.")
		_, err = conn.WriteToUDP(response, clietnAddr)
		if err != nil {
			fmt.Println("Send response failed: ", err)
		}
	}
}
