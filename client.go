package main

import (
	"RandomString"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Write a complete message
func testWriteWhole(client *net.TCPConn, buf *bytes.Buffer) error {
	if _, err := client.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
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
	return nil
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
	return nil
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
	return nil
}

// Write 3 message in one package
func testMutipleMessage(client *net.TCPConn, buf *bytes.Buffer) error {
	data := bytes.Repeat(buf.Bytes(), 3)
	if _, err := client.Write(data); err != nil {
		return err
	}
	return nil
}

func RandomStrings(limit int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	str_len := r.Intn(limit)
	return RandomString.RandomString(uint(str_len))
}

func RunTest(host string, port int) error {
	addr := net.TCPAddr{
		IP:   net.ParseIP(host),
		Port: port,
	}

	client, err := net.DialTCP("tcp", nil, &addr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	message := RandomStrings(100)

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, uint16(len(message))); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, []byte(message)); err != nil {
		return err
	}

	fmt.Println(string(buf.Bytes()))

	if err := testSplitBoth(client, buf); err != nil {
		return err
	}
	if i := 1; i == 0 {
		if err := testWriteWhole(client, buf); err != nil {
			return err
		}
		if err := testSplitHeader(client, buf); err != nil {
			return err
		}

		if err := testSplitBody(client, buf); err != nil {
			return err
		}

		if err := testSplitBoth(client, buf); err != nil {
			return err
		}
		if err := testMutipleMessage(client, buf); err != nil {
			return err
		}
	}
	go Handle(client, message)
	return nil
}

func Handle(conn net.Conn, message string) {
	head := make([]byte, 2)
	defer conn.Close()
	for {
		head_len, err := io.ReadFull(conn, head)
		if head_len != 0 {
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
			//if bytes.Equal([]byte(message), body) {
			fmt.Println("len", info_len, body_len, string(body))
		}
		//}
	}
}

func main() {
	defaultHost := "127.0.0.1"
	defaultPort := 10086
	if len(os.Args) > 1 {
		if port, err := strconv.Atoi(os.Args[1]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			defaultPort = port
		}
		if len(os.Args) > 2 {
			defaultHost = strings.TrimSpace(os.Args[2])
		}
	} else {
		progname := filepath.Base(os.Args[0])
		fmt.Printf("Usage: %s [port] [host]\n", progname)
		//os.Exit(1)
	}
	for i := 0; i < 1; i++ {
		err := RunTest(defaultHost, defaultPort)
		if err != nil {
			fmt.Println(err)
		}
	}
	for {

	}
}
