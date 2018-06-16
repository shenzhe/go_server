package main

import (
	"comm"
	"config"
	"fmt"
	"io"
	"miop"
	"net"
	"time"
	//"strconv"
)

func main() {
	fmt.Printf("pls enter uid@product:\n")
	var uid string
	fmt.Scan(&uid)
	var lid string
	fmt.Printf("pls end last msgid:\n")
	fmt.Scan(&lid)

	room := getRoom()
	if "" == room {
		fmt.Printf("no active room\n")
		return
	}
	conn, err := net.Dial("tcp", room)
	comm.CheckError(err, "connect error")

	var buf [1024]byte

	msg := miop.MiopDataPack(uint32(123), "i am message")
	b := miop.MiopPack("token:7a99ec51ee732b0536e1f93e0bde7737\nu:"+uid+"\nlid:"+lid, msg)
	fmt.Printf("send len: %d\n", len(b))
	_, errr := conn.Write(b)
	comm.CheckError(errr, "send error")
	defer conn.Close()
	i := 1
	for {
		rlen, err := conn.Read(buf[0:])
		if nil != err {
			if err == io.EOF {
				println("receive end")
				break
			} else {
				println("server close\n")
			}
		}
		datas, _ := miop.MiopUnpack(buf[:rlen])
		for _, data := range datas {
			msgs := data.ParseData()
			for _, rmsg := range msgs {
				fmt.Printf("receieve: id:%d msg:%s\n", rmsg.Id, rmsg.Msg)
			}
		}
		//fmt.Printf("receieve: %s \n", buf)
		//i++
		//says := uid + "say: " + strconv.Itoa(i)
		//conn.Write([]byte(says))
		time.Sleep(time.Duration(i) * time.Second)
	}

}

func getRoom() string {
	conn, err := net.Dial("tcp", config.DispatcherHost)
	comm.CheckError(err, "connect err")
	defer conn.Close()

	var buf [100]byte
	rlen, err := conn.Read(buf[0:])

	if nil != err {
		if err != io.EOF {
			comm.CheckError(err, "read error")
		}
	}

	return string(buf[0:rlen])

}
