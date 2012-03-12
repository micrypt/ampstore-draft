// Ampstore.go : Draft version of the Ampify data store
package main

import (
	"fmt"
	"net"
	"os"
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

const SOCK = "/tmp/ampstore.sock"

// This is spawned into a light thread that then handles following requests to the connection. Needs to loop and handle each query.
func (s *Server) handleConn(c net.Conn) {
	defer c.Close()
	buf := make([]byte, 512)
	for {
		nr, err := c.Read(buf)
		fmt.Println("nr: ", nr)
		if err != nil {
			panic(fmt.Sprintf("Read: %v", err))
		}
		fmt.Println("buf: ", buf)
		nw, err := c.Write(buf)
        key := "test"
        s.store.Set(&key, &buf) // Testing Set
		fmt.Println("nw: ", nw)
		if err != nil {
			panic(fmt.Sprintf("Write: %v", err))
		}
	}
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
