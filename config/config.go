package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Addr struct {
		ServerAddr string
		RpcAddr    string
	}
	MaxConn int
}

func LoadConfig(filename string) (*Config, error) {
	byte, err := ioutil.ReadFile(filename)
	if err == nil {
		cfg := Config{}
		err = json.Unmarshal(byte, &cfg)
		if err == nil {
			return &cfg, err
		}
	}
	fmt.Println("AAA")
	return nil, err
}
