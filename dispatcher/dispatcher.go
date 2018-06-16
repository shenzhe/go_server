package dispatcher

import (
	"net"
	//    "fmt"
	"comm"
	"config"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"miop"
	"strconv"
	"time"
)

type RoomInfo struct {
	Rid string `json:"rid"`
	Ol  int    `json:"ol"`
}

type Dispatcher struct {
	RoomList []RoomInfo `json:"rooms"`
}

var Dis Dispatcher

func LoadRoom() {
	roomBuf, err := ioutil.ReadFile(config.RLFile)
	comm.CheckError(err, "read file "+config.RLFile+" error")
	//fmt.Printf("config: %s\n", roomBuf)

	err = json.Unmarshal(roomBuf, &Dis)
	comm.CheckError(err, "json decode err")
	//fmt.Printf("decode: %s\n", Dis)
}

func (dis *Dispatcher) GetAllStats() {
	for index, roomInfo := range dis.RoomList {
		ol, err := GetRoomStats(roomInfo.Rid)
		if nil != err {
			comm.Log("get romm:%s err:%s", roomInfo.Rid, err.Error())
		}
		//fmt.Printf("index: %d, rid: %s, ol: %d", index, roomInfo.Rid, ol)
		dis.RoomList[index].Ol = ol
	}
}

func (dis *Dispatcher) GetRoom() string {
	sum := 0

	for _, roomInfo := range dis.RoomList {
		//fmt.Printf("rid: %s ol: %d",  roomInfo.Rid, roomInfo.Ol)
		if roomInfo.Ol < 0 {
			continue
		}
		sum += roomInfo.Ol
	}

	if sum < 1 {
		return ""
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rand := r.Intn(sum)
	s := 0
	//fmt.Printf("rand:%d, sum:%d\n", rand, sum)

	for _, roomInfo := range dis.RoomList {
		if roomInfo.Ol < 0 {
			continue
		}
		s += roomInfo.Ol
		if rand < s {
			//fmt.Printf("success s: %d, rand: %d\n", s, rand)
			return roomInfo.Rid
		}
	}

	return ""
}

func GetRoomStats(rid string) (int, error) {
	conn, err := net.Dial("tcp", rid)
	if nil != err {
		return -1, err
	}
	defer conn.Close()

	b := miop.MiopPack("token:"+config.AdminToken+"\nu:"+config.AdminUid, []byte(""))
	//fmt.Printf("send len: %d\n", len(b))
	_, errr := conn.Write(b)
	if nil != errr {
		return -1, err
	}
	buf := make([]byte, config.BuffLen)
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
			prop := data.ParseProp()
			ol := prop.GetPropKey("ol")
			if "" == ol {
				continue
			}
			return strconv.Atoi(ol)
		}
	}

	return -1, nil

}
