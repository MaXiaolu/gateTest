package main

import (
	"github.com/MaXiaolu/gateTest/game"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	arith := new(gameserver.Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
	for {

	}
}
