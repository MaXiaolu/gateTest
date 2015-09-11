package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

type Server struct {
	conn net.Conn
}

var (
	conns       map[int]net.Conn = make(map[int]net.Conn)
	next_serial                  = 1
)

func GetNextSerial() int {
	serial := next_serial
	for {
		if conns[serial] != nil {
			serial += 1
		} else {
			break
		}
	}
	next_serial = serial + 1
	return serial
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:10086")
	checkErr(err)
	listen, err := net.ListenTCP("tcp", addr)
	checkErr(err)
	fmt.Println("Start server...")
	for {
		conn, err := listen.Accept()
		checkErr(err)
		server := new(Server)
		server.InitServer(conn)
		go server.Handle() // 每次建立一个连接就放到单独的线程内做处理
	}
}

func (s *Server) InitServer(conn net.Conn) {
	serial := GetNextSerial()
	conns[serial] = conn
	s.conn = conn
	fmt.Println("AAAA")
}

func (s *Server) Handle() {
	conn := s.conn
	head := make([]byte, 2)
	defer conn.Close()
	for {
		_, err := io.ReadFull(conn, head)
		if err != nil {
			fmt.Println(err)
			break
		}
		var info_len uint16
		buf := bytes.NewReader(head)
		if err := binary.Read(buf, binary.BigEndian, &info_len); err != nil {
			fmt.Println("binary.Read failed:", err)
			break
		}
		body := make([]byte, info_len)
		body_len, err := io.ReadFull(conn, body)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("len", info_len, body_len, string(body))
		WriteTo(1, body)
	}
}

func WriteTo(serial int, data []byte) error {
	conn := conns[serial]
	if conn != nil {
		buf := new(bytes.Buffer)
		if err := binary.Write(buf, binary.BigEndian, uint16(len(data))); err != nil {
			return err
		}
		if err := binary.Write(buf, binary.LittleEndian, data); err != nil {
			return err
		}
		conn.Write(buf.Bytes())
	}
	return nil
}

func IsConnClosed(serial int) {
	conn := conns[serial]
	if conn != nil {
		if _, err := conn.Write([]byte("")); err != nil {
			conns[serial] = nil
		}
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
