package user

import (
	"bytes"
	"comm"
	"config"
	"entity"
	"fmt"
	"miop"
	"net"
	"strconv"
	"time"
)

var (
	userlist  = make([]User, 1000)
	ConnCount = 0
)

type User struct {
	Conn     net.Conn
	C        chan []byte
	Uid      string
	More     []byte
	LastTime int64
	Ue       *entity.UserEntity
}

var UserList = make(map[string]*User)

func CheckUserExits(uid string) bool {
	_, ok := UserList[uid]
	if ok {
		return true
	}
	return false
}

func AddUser(u *User) *User {
	if CheckUserExits(u.Uid) {
		return nil
	}

	comm.Log("user %s add", u.Uid)

	//save to db

	user := entity.GetUser(u.Uid)

	if nil == user {
		now := time.Now().Unix()
		entity.AddUser(u.Uid, u.Conn.RemoteAddr().String(), 1000, now, now, now)
	}

	return SetUser(u)
}

func SetUser(u *User) *User {
	UserList[u.Uid] = u
	return u
}

func DelUser(uid string) bool {
	if CheckUserExits(uid) {
		delete(UserList, uid)
		return true
	}
	return false
}

func GetUser(uid string) *User {
	u, ok := UserList[uid]
	if ok {
		return u
	}
	return nil
}

func (u *User) PutChan() bool {
	buf := make([]byte, config.BuffLen)
	rlen, err := u.Conn.Read(buf)
	fmt.Printf("read len: %d", rlen)
	if nil == err {
		comm.Debug("conn read %x", buf[:rlen])
		if len(buf) > 0 {
			u.C <- buf[:rlen]
		}
		u.UpTime()
		return u.OnRecieve()
	} else {
		comm.Log("close connection: %s", u.Conn.RemoteAddr().String())
		u.Conn.Close()
		DelUser(u.Uid)
		return false
	}
}

func (u *User) OnRecieve() bool {
	select {
	case msg := <-u.C:
		return u.parseMsg(msg)
	default:
		return u.CheckHeartBeat()

	}
}

func (u *User) CheckHeartBeat() bool {
	comm.Log("hb check %s", u.Uid)
	return u.LastTime+config.HeartbeatTime > time.Now().Unix()
}

func (u *User) UpTime() {
	u.LastTime = time.Now().Unix()
}

func (u *User) Send(msg []byte) {
	wrote, err := u.Conn.Write(msg)
	comm.CheckError(err, "write: "+string(wrote)+" bytes")
	comm.Log("send to %s : %x", u.Uid, msg)
}

func (u *User) SendOfflineMsg(id int) {
	msgs := comm.GetMsgList(id)
	if len(msgs) > 0 {
		comm.Log("uid:%s get %d offline msg", u.Uid, len(msgs))
		msgs = miop.ParseMsgList(msgs)
		msg := miop.MiopPack("", bytes.Join(msgs, []byte("")))
		u.Send(msg)
	}
	return
}

func (u *User) parseMsg(msg []byte) bool {
	comm.Log("parsemsg %x len %d", msg, len(msg))
	if len(u.More) > 0 {
		//msg = comm.MergeByte(u.More, msg)
		msg = append(u.More, msg...)
	}
	datas, more := miop.MiopUnpack(msg)
	if len(more) > 0 {
		u.More = more
	}
	comm.Log("unpack %d", len(datas))
	result := true
	for _, data := range datas {
		if "" == u.Uid { //first connect
			prop := data.ParseProp()
			if ok := prop.CheckToken(); !ok {
				comm.Log("check token fail")
				result = false
				continue
			}
			uname, p := prop.GetUser()
			if checkAdmin(uname) {
				token := prop.GetPropKey("token")
				if token == config.AdminToken {
					u.Send(miop.MiopPack("ol:"+strconv.Itoa(ConnCount), []byte("")))
				}
				return false
			}
			if "" != uname {
				u.Uid = uname + ":" + p
				AddUser(u)
				lastId := prop.GetLastId()
				u.SendOfflineMsg(lastId)
				result = true
				continue
			}
		} else {
			u.Send(msg)
		}
	}

	return result
}

func checkAdmin(u string) bool {
	return u == config.AdminUid
}

func BoardCast(msg []byte) {
	for _, u := range UserList {
		comm.Log("send to %s", u.Uid)
		u.parseMsg(msg)
	}
}

func GetConnCount() int {
	return ConnCount
}
