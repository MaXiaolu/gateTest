package main

import (
	"github.com/MaXiaolu/gateTest/config"
	"github.com/MaXiaolu/gateTest/server"
	"log"
	"os"
)

func main() {
	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	s := server.NewServer()
	err = s.Start(cfg.Addr.ServerAddr, cfg.Addr.RpcAddr)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}
