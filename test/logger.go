package main

import (
    "comm"
)


func main () {
    comm.GetLogger(comm.DEBUG).Println("i am debug log")
    comm.GetLogger(comm.LOG).Println("i am push log")
    comm.GetLogger(comm.DEBUG).Println("i am debug log2")
    comm.GetLogger(comm.LOG).Println("i am push log2")
    comm.GetLogger(comm.DEBUG).Println("i am debug log3")
    comm.GetLogger(comm.LOG).Println("i am push log3")
    comm.Debug("i am %s", "string")
    comm.Debug("i am string too")
}
