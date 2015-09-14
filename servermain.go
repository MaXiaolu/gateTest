package main

import (
	"fmt"
	"github.com/MaXiaolu/gateTest/gameserver"
	"github.com/MaXiaolu/gateTest/server"
	"log"
	"net/rpc"
)

func main() {
	ConnToGame()
	s := server.NewServer()
	s.Start()
}

func ConnToGame() {
	client, err := rpc.DialHTTP("tcp", ":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	args := &gameserver.Args{17, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Println("Arith: %d * %d = %d", args.A, args.B, reply)

	// Asynchronous call
	quotient := new(gameserver.Quotient)
	divCall := client.Go("Arith.Divide", args, quotient, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	// check errors, print, etc.
	fmt.Println(replyCall.Reply)
}
