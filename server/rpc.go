package server

import (
	"fmt"
	//"github.com/MaXiaolu/gateTest/config"
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

func StarRpc(RpcAddr string) error {
	log.Printf("gamed started ")
	for {
		l, err := net.Listen("tcp", RpcAddr)
		if err != nil {
			return err
		}
		arith := new(Arith)
		rpc.Register(arith)
		rpc.HandleHTTP()
		http.Serve(l, nil)
	}
}

func CallRpc(RpcAddr, method string, A, B int) int {
	reply := 0
	client, err := rpc.DialHTTP("tcp", RpcAddr)
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
