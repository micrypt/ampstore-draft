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

const BUFLEN = 12
const SOCK = "/tmp/ampstore.sock"

func (s *Server) handleConn(c net.Conn) {
	defer c.Close()
	buf := make([]byte, BUFLEN)
    cmdVal := []byte{}
    var cmd string
    var params []string

	for {
		nr, err := c.Read(buf)
		fmt.Println("nr: ", nr)
		if err != nil {
			panic(fmt.Sprintf("Read: %v", err))
		}
        if nr > 0 {
          cmdVal = append(cmdVal, buf[0:nr]...)
          cmd, params, _, _ = parseRequest(cmdVal)
        }
        fmt.Printf("cmd: %v, params: %v\n", cmd, params)
        fmt.Printf("cmdVal: %v\n", cmdVal)
		nw, err := c.Write(cmdVal)
		fmt.Println("nw: ", nw)
        key := "test" // Dummy key
        s.store.Set(&key, &cmdVal) // Testing Set
		if err != nil {
			panic(fmt.Sprintf("Write: %v", err))
		}
	}
}

func parseRequest(b []byte) (cmd string, params, vals []string, err error) {
  cmd = string(b[0])
  paramsLen, err := strconv.ParseInt(string(b[1]), 10, 32)
  params = make([]string, paramsLen)
  var i int64
  for i = 0; i < paramsLen; i++ {
    params[i] = string(b[2])
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
