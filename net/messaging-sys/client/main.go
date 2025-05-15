package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("连接服务器失败")
		return
	}

	go client.ProcessResponse()
	fmt.Println("连接服务器成功")

	client.Run()
}
