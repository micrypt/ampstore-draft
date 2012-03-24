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
	store *KVStore
	sock  string
	stype int
}

const (
	Nil = iota
	TCP_SOCK
	UNIX_SOCK
)

const (
	BUFLEN       = 512
	PARAMOFFSET  = 2
	ERROR_STRING = "ERROR"
	OK_STRING    = "OK"
)

var COMMANDS_MAP = map[rune]func(*KVStore, []string, *[]byte) error{
	'\x01': getCommand,
	'\x02': setCommand,
}

func getCommand(s *KVStore, values []string, resp *[]byte) error {
	err := s.Get(values[0], resp)
	return err
}

func setCommand(s *KVStore, values []string, resp *[]byte) error {
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
		if err != nil {
			panic(fmt.Sprintf("Read error: %v", err))
		}
		if nr > 0 {
			cmdSlice = append(cmdSlice, buf[0:nr]...)
			cmd, params, done, _ := parseRequest(cmdSlice)
			if done && params != nil {
				err1 := COMMANDS_MAP[cmd](s.store, params, &respSlice)
				if err1 != nil {
					c.Write([]byte(ERROR_STRING))
					continue
				}
				_, err := c.Write(respSlice)
				if err != nil {
					panic(fmt.Sprintf("Write error: %v", err))
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
		rStart := offset + 1
		rEnd := rStart + length
		params[i] = string(b[rStart:rEnd])
		offset = rEnd
	}
	totalLen = PARAMOFFSET + paramsLen + totalLen
	if totalLen == int(len(b)) {
		done = true
	}
	return
}

func (s *Server) Init() {
	var listener net.Listener
	var err error
	if s.stype == UNIX_SOCK {
		os.Remove(s.sock)
		defer os.Remove(s.sock)
		listener, err = net.Listen("unix", s.sock)
		if err != nil {
			panic(fmt.Sprintf("Error: %v\n", err))
		}
	} else {
		listener, err = net.Listen("tcp", s.sock)
		if err != nil {
			panic(fmt.Sprintf("Error: %v\n", err))
		}
	}
	fmt.Println("Listening on", listener.Addr())
	for {
		c, err := listener.Accept()
		if err != nil {
			panic(fmt.Sprintf("Accept: %v", err))
		}
		go s.handleConn(c)
	}
}

func runServer(socket string, stype int) {
	server := &Server{NewKVStore(), socket, stype}
	server.Init()
}
