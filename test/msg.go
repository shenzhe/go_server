package main

import (
    "fmt"
    "comm"
    "miop"
)

func main() {
    ok, id := comm.AddMsg("i am baordmsg")
    if ok {
        fmt.Printf("add msg succ :%d\n", id)
    }

    msgs := comm.GetMsgList(0)

    msgs = miop.ParseMsgList(msgs)

    for _, msg := range msgs {
        fmt.Printf("get msg: %s\n", msg)
    }
}
