package main

import (
	"comm"
	"config"
	"server"
	"signal"
)

func main() {
	config.Init()
	go signal.Start()
	go server.BoardCast()
	ln := server.Start(config.HostAndPort)
	for {
		conn, err := ln.Accept()
		comm.CheckError(err, "Accept error")
		go server.HandlerConn(conn)
	}
}
