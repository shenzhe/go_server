package signal

import (
    "os"
    osignal "os/signal"
    "syscall"
    "server"
    "comm"
)

type sigHandler func (s os.Signal, arg interface{})

type sigSet struct {
    m map[os.Signal]sigHandler
}

func sigInit() *sigSet {
    sig := new(sigSet)
    sig.m = make(map[os.Signal]sigHandler)
    return sig
}

func (ss *sigSet) add(s os.Signal, handler sigHandler) {
    if _, ok := ss.m[s]; !ok {
        comm.Log("add signal %d", s)
        ss.m[s] = handler
    }
}

func (ss *sigSet) handler(s os.Signal, arg interface{}) bool {
    if _, ok := ss.m[s]; ok {
        ss.m[s](s, arg)
        return true
    }
    
    return false
}

func Start() {
    comm.Log("signal start")
    ss := sigInit()
    comm.Log("signal init")
    handler := func(s os.Signal, arg interface{}){
                   comm.Log("receive sig %d", s)
                   server.Stop()
    } 
    ss.add(syscall.SIGHUP, handler)
    
    for {
        c := make(chan os.Signal)
        var sigs []os.Signal
        for sig := range ss.m {
            sigs = append(sigs, sig)
        }
        
        osignal.Notify(c)
        
        sig := <-c
        
        ss.handler(sig, nil)
    }
}

