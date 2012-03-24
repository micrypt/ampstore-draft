package main

import "flag"

func main() {
	var stype int
	modes := map[string]func(socket string, stype int){"client": runClient, "server": runServer}
	mode := flag.String("mode", "", "What mode to start in")
	socket := flag.String("socket", "", "Unix domain socket to connect to")
	stype_raw := flag.String("stype", "", "Socket type to use")
	flag.Parse()
	if *stype_raw == "unix" {
		stype = UNIX_SOCK
	} else {
		stype = TCP_SOCK
	}
	modes[*mode](*socket, stype)
}
