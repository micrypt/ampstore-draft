// Ampstore.go : Draft version of the Ampify data store
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	conn      net.Conn
	sock_mode int
	socket    string
}

var COMMANDS = map[string]rune{
	"GET":  '\x01',
	"SET":  '\x02',
	"SCAN": '\x03',
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
	_, err := fmt.Fprint(cl.conn, msg)
	fmt.Println("Sending: ", msg)
	if err != nil {
		panic(fmt.Sprintf("Transmission error: %v", err))
	}
	_, err1 = cl.conn.Read(b[0:])
	if err1 != nil {
		panic(fmt.Sprintf("Receive error: %v", err1))
	}
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
		fmt.Print(">> ") // Show prompt
		read, err := buf.ReadString('\n')
		if err != nil {
			fmt.Println("")
			break
		}
		if len(read) == 0 {
			continue
		}
		read = strings.Trim(read, "\n")
		sects := strings.Split(read, " ")
		paramsLen := len(sects) - 1
		cmd := string([]rune{COMMANDS[strings.ToUpper(sects[0])]})
		if paramsLen > 0 {
			cmd = cmd + fmt.Sprint(paramsLen)
			for _, v := range sects[1:] {
				cmd = cmd + fmt.Sprint(len([]byte(v))) + v
			}
		}
		resp, err := client.SendCommand(cmd)
        // TODO: Only read required bytes of response
		if err != nil {
			panic(fmt.Sprintf("Response: %v, Error: %v", resp, err))
		}
		fmt.Println("Response: ", resp)
	}
}
