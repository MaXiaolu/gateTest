package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Args struct {
	A, B int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func StarRpc() {
	log.Printf("gamed started ")
	for {
		l, e := net.Listen("tcp", ":1234")
		if e != nil {
			log.Fatal("listen error:", e)
		}
		arith := new(Arith)
		rpc.Register(arith)
		rpc.HandleHTTP()
		http.Serve(l, nil)
	}
}

func CallRpc(method string, A, B int) int {
	reply := 0
	client, err := rpc.DialHTTP("tcp", ":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	} else {
		fmt.Println("A_B :", A, B)
		args := &Args{A, B}
		if err := client.Call(method, args, &reply); err != nil {
			log.Fatal("arith error:", err)
		}
		client.Close()
	}
	return reply
}
