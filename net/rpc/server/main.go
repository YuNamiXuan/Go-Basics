package main

import (
	"errors"
	"log"
	"net"
	"net/rpc"
)

type MathService struct{}

type Args struct {
	A, B float64
	Op   string
}

type Reply struct {
	Ans float64
}

func (m *MathService) Calculate(args *Args, reply *Reply) error {
	log.Printf("Received request: %f %s %f", args.A, args.Op, args.B)
	switch args.Op {
	case "+":
		reply.Ans = args.A + args.B
	case "-":
		reply.Ans = args.A - args.B
	case "*":
		reply.Ans = args.A * args.B
	case "/":
		if args.B == 0 {
			return errors.New("division by zero")
		}
		reply.Ans = args.A / args.B
	default:
		return errors.New("invalid operator")
	}
	return nil
}

func main() {
	mathService := new(MathService)
	rpc.Register(mathService)

	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal("Listen error: ", err)
	}
	log.Println("RPC Server running on: 8080...")
	rpc.Accept(listener)
}
