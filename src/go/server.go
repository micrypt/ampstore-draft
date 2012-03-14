// Ampstore.go : Draft version of the Ampify data store
package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
)

type Server struct {
	store     *KVStore
	Addr      string
	sock_mode int
}

const (
	Nil = iota
	TCP_SOCK
	UNIX_SOCK
)

const (
	BUFLEN       = 12
	PARAMOFFSET  = 2
	SOCK         = "/tmp/ampstore.sock"
	ERROR_STRING = "ERROR"
	OK_STRING = "OK"
)

var COMMANDS_MAP = map[rune]func(*KVStore, []string, *[]byte) error{
	'\x01': getCommand,
	'\x02': setCommand,
}

func getCommand(s *KVStore, values []string, resp *[]byte) error {
    fmt.Println("Running GET with: ", values)
	err := s.Get(values[0], resp)
	return err
}

func setCommand(s *KVStore, values []string, resp *[]byte) error {
    fmt.Println("Running SET with: ", values)
    key := values[0]
    var val []byte
    for _, v := range values[1:] {
      val = append(val, []byte(v)...)
    }
    err := s.Set(key, &val, resp)
	return err
}

func (s *Server) handleConn(c net.Conn) {
	defer c.Close()
	buf := make([]byte, BUFLEN)
	var cmdSlice, respSlice []byte
	for {
		nr, err := c.Read(buf)
		fmt.Println("nr: ", nr)
		if err != nil {
			panic(fmt.Sprintf("Read: %v", err))
		}
		if nr > 0 {
			cmdSlice = append(cmdSlice, buf[0:nr]...)
			cmd, params, done, _ := parseRequest(cmdSlice)
			if done && params != nil {
				fmt.Printf("cmd: %v, params: %v\n", cmd, params)
				fmt.Printf("cmdSlice: %v\n", cmdSlice)
				err1 := COMMANDS_MAP[cmd](s.store, params, &respSlice)
				if err1 != nil {
					c.Write([]byte(ERROR_STRING))
					continue
				}
                fmt.Println("respSlice: ", respSlice)
				nw, err := c.Write(respSlice)
				fmt.Println("nw: ", nw)
				if err != nil {
					panic(fmt.Sprintf("Write: %v", err))
				}
				cmdSlice = []byte{}
			}
		}
	}
}

func makeInt(b byte) (i int, err error) {
	i, err = strconv.Atoi(string(b))
	return
}

func parseRequest(b []byte) (cmd rune, params []string, done bool, err error) {
	cmd = bytes.Runes(b[:1])[0]
	paramsLen, err := makeInt(b[1])
	params = make([]string, paramsLen)
	var offset, length, totalLen, i int
	offset = PARAMOFFSET
	for i = 0; i < paramsLen; i++ {
		length, _ = makeInt(b[offset])
		totalLen = totalLen + length
		fmt.Println("Param length: ", length)
		rStart := offset + 1
		rEnd := rStart + length
		fmt.Println("rStart: ", rStart, " rEnd: ", rEnd)
		params[i] = string(b[rStart:rEnd])
		fmt.Println("Current param: ", params[i])
		offset = rEnd
	}
	totalLen = PARAMOFFSET + paramsLen + totalLen
	if totalLen == int(len(b)) {
		done = true
	}
	return
}

func (s *Server) Init() {
	if s.sock_mode == UNIX_SOCK {
		fmt.Println("Address: ", s.Addr)
		os.Remove(s.Addr)
		defer os.Remove(s.Addr)
		listener, err := net.Listen("unix", s.Addr)
		if err != nil {
			panic(fmt.Sprintf("Error: %v\n", err))
		}
		fmt.Println("Listening on", listener.Addr())
		for {
			c, err := listener.Accept()
			if err != nil {
				panic(fmt.Sprintf("Accept: %v", err))
			}
			fmt.Println("Handling connection")
			go s.handleConn(c)
		}
	}
}

func NewServer() *Server {
	server := &Server{NewKVStore(), SOCK, UNIX_SOCK}
	return server
}

func runServer(socket string) {
	server := NewServer()
	server.Addr = socket
	server.Init()
}
