package main

import "flag"

func main() {
	modes := map[string]func(socket string){"client": runClient, "server": runServer}
	mode := flag.String("mode", "", "What mode to start in")
	socket := flag.String("socket", "", "Unix domain socket to connect to")
	flag.Parse()
	modes[*mode](*socket)
}
