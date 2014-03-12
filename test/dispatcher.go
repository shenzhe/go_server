package main

import (
    "net"
    "fmt"
    "comm"
    "config"
    "io"
)


func main() {
    conn, err := net.Dial("tcp", config.DispatcherHost)
    comm.CheckError(err, "connect err")
    defer conn.Close()
    
    var buf [100]byte
    rlen, err := conn.Read(buf[0:])
   
    fmt.Printf("receive len: %d", rlen) 
        if nil != err {
            if err != io.EOF {
                comm.CheckError(err, "read error")
            }
        }

    
    fmt.Printf("get rid %s", buf[0:rlen])    


}
