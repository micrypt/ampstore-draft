// Ampstore.go : Draft version of the Ampify data store
package main

import (
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
	BUFLEN      = 12
	PARAMOFFSET = 2
	SOCK        = "/tmp/ampstore.sock"
)

func (s *Server) handleConn(c net.Conn) {
	defer c.Close()
	buf := make([]byte, BUFLEN)
	cmdSlice := []byte{}
	for {
		nr, err := c.Read(buf)
		fmt.Println("nr: ", nr)
		if err != nil {
			panic(fmt.Sprintf("Read: %v", err))
		}
		if nr > 0 {
			cmdSlice = append(cmdSlice, buf[0:nr]...)
			cmd, params, _, done, _ := parseRequest(cmdSlice)
			if done && cmd != "" && params != nil {
				fmt.Printf("cmd: %v, params: %v\n", cmd, params)
				fmt.Printf("cmdSlice: %v\n", cmdSlice)
				nw, err := c.Write(cmdSlice)
				fmt.Println("nw: ", nw)
				key := "test"              // Dummy key
				s.store.Set(&key, &cmdSlice) // Testing Set
				if err != nil {
					panic(fmt.Sprintf("Write: %v", err))
				}
                cmdSlice = []byte{}
			}
		}
	}
}

func makeInt(b byte) (i int64, err error) {
	i, err = strconv.ParseInt(string(b), 10, 32)
	return
}

func parseRequest(b []byte) (cmd string, params, vals []string, done bool, err error) {
    //bufferLen := len(b)
	cmd = string(b[0])
	paramsLen, err := makeInt(b[1])
	params = make([]string, paramsLen)
	var offset, length, i int64
	offset = PARAMOFFSET
	for i = 0; i < paramsLen; i++ {
		length, _ = makeInt(b[offset])
        fmt.Println("Param length: ", length)
        rStart := offset+1
        rEnd := rStart+length
        fmt.Println("rStart: ", rStart, " rEnd: ", rEnd)
		params[i] = string(b[rStart:rEnd])
        fmt.Println("Current param: ", params[i])
		offset = rEnd
	}
    done = true
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
