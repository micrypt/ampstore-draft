// Ampstore.go : Draft version of the Ampify data store
package main

import (
    "bufio"
	"fmt"
	"net"
    "os"
)

type Client struct {
	conn      net.Conn
	sock_mode int
	socket    string
}

func (cl *Client) Connect() error {
	if cl.sock_mode == UNIX_SOCK {
		// Create the connection
		conn, err := net.Dial("unix", cl.socket)
        if err != nil {
          panic(fmt.Sprintf("Error creating connection: %v", err))
        }
		cl.conn = conn
		// Print info about the unix socket
		fmt.Println("Connected to", conn.RemoteAddr())
	}
	return nil
}

func (cl *Client) SendCommand(msg string) (b [512]byte, err1 error) {
	//_, err := cl.conn.Write(msg)
	fmt.Fprint(cl.conn, msg)
    fmt.Println("Sending: %v", msg)
//	if err != nil {
//		panic(fmt.Sprintf("Transmission error: %v", err))
//	}
	//_, err1 = cl.conn.Read(b[0:])
	//if err1 != nil {
	//	panic(fmt.Sprintf("Receive error: %v", err1))
	//}
	return
}

func NewClient() *Client {
	cl := &Client{nil, UNIX_SOCK, SOCK}
	return cl
}

func runClient(socket string) {
	client := NewClient()
	client.socket = socket
	client.Connect()
    buf := bufio.NewReader(os.Stdin)
	for {
	    fmt.Print(">>") // Show prompt
        read, err := buf.ReadString('\n')
        if err != nil {
          fmt.Println("")
          break
        }
        line := read[0:len(read)-1]
        if len(line) == 0 {
          continue
        }
		if resp, err := client.SendCommand(line); err != nil {
          panic(fmt.Sprintf("Response: %v, Error: %v", resp, err))
		}
	}
}
