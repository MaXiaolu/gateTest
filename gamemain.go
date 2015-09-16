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
	if err = server.StarRpc(cfg.Addr.RpcAddr); err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}
