package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strings"
)

type Args struct {
	A, B float64
	Op   string
}

type Reply struct {
	Ans float64
}

func main() {
	client, err := rpc.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Dial error: ", err)
		return
	}
	defer client.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		var args Args
		var reply Reply
		fmt.Print("Enter calculation (e.g., '10 20 +'): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		_, err := fmt.Sscanf(input, "%f %f %s", &args.A, &args.B, &args.Op)
		if err != nil {
			fmt.Println("Invalid input format. Example: '10 20 +'")
			continue
		}
		err = client.Call("MathService.Calculate", &args, &reply)
		if err != nil {
			fmt.Println("RPC Call error: ", err)
			continue
		}
		fmt.Println(reply.Ans)
	}
}
