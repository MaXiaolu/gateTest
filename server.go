package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

type Server struct {
	conns       map[int]net.Conn
	next_serial int
}

func NewServer() *Server {
	return &Server{conns: make(map[int]net.Conn), next_serial: 1}
}

func (s *Server) Start() {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:10086")
	checkErr(err)
	listen, err := net.ListenTCP("tcp", addr)
	checkErr(err)
	fmt.Println("Start server...")
	for {
		conn, err := listen.Accept()
		checkErr(err)
		serial := s.GetNextSerial()
		s.conns[serial] = conn
		go s.Handle(serial) // 每次建立一个连接就放到单独的线程内做处理
	}
}

func (s *Server) Handle(serial int) {
	conn := s.conns[serial]
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
		//s.WriteTo(1, body)
	}
}

func (s *Server) GetNextSerial() int {
	serial := s.next_serial
	for {
		if s.conns[serial] != nil {
			serial += 1
		} else {
			break
		}
	}
	s.next_serial = serial + 1
	return serial
}

func (s *Server) WriteTo(serial int, data []byte) error {
	conn := s.conns[serial]
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

func (s *Server) Close(serial int) {
	conn := s.conns[serial]
	if conn != nil {
		if _, err := conn.Write([]byte("")); err != nil {
			s.conns[serial] = nil
		}
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
