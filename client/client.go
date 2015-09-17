package main

import (
	"RandomString"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/MaXiaolu/gateTest/config"
	"io"
	"log"
	"math/rand"
	"net"
	"time"
)

// Write a complete message
func testWriteWhole(client *net.TCPConn, buf *bytes.Buffer) error {
	if _, err := client.Write(buf.Bytes()); err != nil {
		return err
	}
	return EqualRead(client, buf.Bytes())
}

// Write part of header first
func testSplitHeader(client *net.TCPConn, buf *bytes.Buffer) error {
	data := buf.Bytes()
	if _, err := client.Write(data[:1]); err != nil {
		return err
	}
	if _, err := client.Write(data[1:]); err != nil {
		return err
	}
	return EqualRead(client, data)
}

// Write header and part of body first
func testSplitBody(client *net.TCPConn, buf *bytes.Buffer) error {
	data := buf.Bytes()
	split := len(data) / 2
	if _, err := client.Write(data[:split]); err != nil {
		return err
	}
	if _, err := client.Write(data[split:]); err != nil {
		return err
	}
	return EqualRead(client, data)
}

// Write part of header, then left of header and part of body, then left body
func testSplitBoth(client *net.TCPConn, buf *bytes.Buffer) error {
	data := buf.Bytes()
	split := len(data)/2 + 1
	if _, err := client.Write(data[:1]); err != nil {
		return err
	}
	if _, err := client.Write(data[1:split]); err != nil {
		return err
	}
	if _, err := client.Write(data[split:]); err != nil {
		return err
	}
	return EqualRead(client, data)
}

// Write 3 message in one package
func testMutipleMessage(client *net.TCPConn, buf *bytes.Buffer) error {
	data := bytes.Repeat(buf.Bytes(), 3)
	if _, err := client.Write(data); err != nil {
		return err
	}
	return EqualRead(client, data)
}

func RandomStrings(limit int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	str_len := r.Intn(limit)
	return RandomString.RandomString(uint(str_len))
}

func RunTest(hostAddress string) {
	addr, err := net.ResolveTCPAddr("tcp", hostAddress)
	if err != nil {
		log.Print(err)
		return
	}

	client, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	message := RandomStrings(100)

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint16(len(message)))
	binary.Write(buf, binary.LittleEndian, []byte(message))

	//fmt.Println(string(buf.Bytes()))

	if err := testSplitBoth(client, buf); err != nil {
		log.Print(err)
	}
	if i := 1; i == 0 {
		if err := testWriteWhole(client, buf); err != nil {
			log.Print(err)
		}
		if err := testSplitHeader(client, buf); err != nil {
			log.Print(err)
		}

		if err := testSplitBody(client, buf); err != nil {
			log.Print(err)
		}

		if err := testSplitBoth(client, buf); err != nil {
			log.Print(err)
		}
		if err := testMutipleMessage(client, buf); err != nil {
			log.Print(err)
		}
	}
}

func EqualRead(conn net.Conn, data []byte) error {
	buf := make([]byte, len(data))
	if body_len, err := io.ReadFull(conn, buf); err != nil {
		return err
	} else {
		fmt.Println(string(data), string(buf))
		if body_len > 0 && !bytes.Equal(data, buf) {
			return errors.New("not equal")
		}
	}
	return nil
}

func main() {
	cfg, err := config.LoadConfig("../config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < cfg.MaxConn; i++ {
		go RunTest(cfg.Addr.ServerAddr)

	}
	for {

	}
}
