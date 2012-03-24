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
		conn, err := net.Dial("unix", cl.socket)
		if err != nil {
			panic(fmt.Sprintf("Error creating connection: %v", err))
		}
		cl.conn = conn
		fmt.Println("Connected to", conn.RemoteAddr())
	} else {
		conn, err := net.Dial("tcp", cl.socket)
		if err != nil {
			panic(fmt.Sprintf("Error creating connection: %v", err))
		}
		cl.conn = conn
		fmt.Println("Connected to", conn.RemoteAddr())
	}
	return nil
}

func (cl *Client) SendCommand(msg string) (b [BUFLEN]byte, err1 error) {
	_, err := fmt.Fprint(cl.conn, msg)
	if err != nil {
		panic(fmt.Sprintf("Transmission error: %v", err))
	}
	_, err1 = cl.conn.Read(b[0:])
	if err1 != nil {
		panic(fmt.Sprintf("Receive error: %v", err1))
	}
	return
}

func runClient(socket string, stype int) {
	client := &Client{nil, stype, socket}
	client.Connect()
	buf := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">> ") // Show prompt
		read, err := buf.ReadString('\n')
		if err != nil {
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
		// TODO: Only read required bytes of response
		resp, err := client.SendCommand(cmd)
		if err != nil {
			panic(fmt.Sprintf("Response: %v, Error: %v", resp, err))
		}
		fmt.Println("Response: ", resp)
	}
}
