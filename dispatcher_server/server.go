package main

import (
	"comm"
	"config"
	"dispatcher"
	"fmt"
	"net"
	"server"
	"time"
)

func main() {
	dispatcher.LoadRoom()
	fmt.Printf("%s", dispatcher.Dis)
	go func() {
		for {
			dispatcher.Dis.GetAllStats()
			time.Sleep(3 * time.Second)
		}
	}()
	ln := server.Start(config.DispatcherHost)
	for {
		conn, err := ln.Accept()
		comm.CheckError(err, "accept error")
		go handlerConn(conn)
	}
}

func handlerConn(conn net.Conn) {
	defer conn.Close()
	rid := dispatcher.Dis.GetRoom()
	fmt.Printf("get rid: %s\n", rid)
	conn.Write([]byte(rid))
}
