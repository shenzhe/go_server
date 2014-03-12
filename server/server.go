package server

import (
    "net"
    "comm"
    "sync"
    "os"
    "time"
    "config"
    "user"
    "entity"
    "miop"
    "bytes"
)

var runing = true
var sending = false
var connCount = 0
var l sync.Mutex

func Stop() {
    runing = false
    comm.Log("server stoping")
    
    if !sending {
        comm.Log("server stop")
        os.Exit(-1)
    }
}

func IsRuning() bool {
    return runing
}

func SendStart() {
    sending = true
    comm.Log("sending....")
}

func SendEnd() {
    sending = false
    comm.Log("send stop....")
    
    if !runing {
        comm.Log("server stop")
        os.Exit(-1)
    }
}

func Start(hostAndPort string) net.Listener {
    listener, err := net.Listen("tcp", hostAndPort)
    comm.CheckError(err, "Listen error")
    comm.Log("Listening to: %s", listener.Addr().String())
    return listener
}

func HandlerConn(conn net.Conn) {
    defer func(conn net.Conn) {
        DelConnCount()
        conn.Close()
    }(conn)

    if !IsRuning() {
        return 
    }

    c := make(chan []byte, config.BuffLen)
    u := &user.User{conn, c, "", nil, time.Now().Unix(), entity.GetEntity()}
    go func() {
        for { 
            time.Sleep(3 * time.Second)
            if !u.CheckHeartBeat() {
                u.Conn.Close()
                return 
            }
        }
        return  
    }()
    AddConnCount()
    for {
        suc := u.PutChan()
        if false == suc {
            break
        }
    }
    
    return
}

func AddConnCount() {
    l.Lock()
    defer l.Unlock()
    connCount++
    user.ConnCount++
}

func DelConnCount() {
    l.Lock()
    defer l.Unlock()
    connCount--
    user.ConnCount--
}

func GetConnCount() int {
    return connCount
}

func BoardCast() {
    for {
        time.Sleep(3 * time.Second)
        sid := comm.GetSendMaxId()
        msg := comm.GetMsgList(sid)
        if nil == msg {
            continue
        }
        SendStart()
        //comm.Log("send boardcast %s", msg)
        sendMsg := miop.ParseMsgList(msg)
        b := miop.MiopPack("", bytes.Join(sendMsg, []byte("")))
        //comm.Log("packmsg %x len %d", b, len(b))
        user.BoardCast(b)
        comm.SetSendMaxId()
        SendEnd()
    }
}

