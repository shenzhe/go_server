package comm

import (
	"bytes"
	"config"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/hoisie/redis"
	"strconv"
)

func CheckError(err error, msg string) {
	if nil != err {
		str := "ERROR: " + msg + " " + err.Error()
		GetLogger(ERROR).Println(msg)
		panic(str)
	}
}

func LogError(err error, msg string) {
	if nil != err {
		GetLogger(ERROR).Println(msg)
	}

}

var client redis.Client

func GetMsg(key string) []byte {
	return GetData(key)
}

func GetMsgList(id int) [][]byte {
	id++
	msgs, err := client.Zrangebyscore(config.BoardCastKey, float64(id), 10000)
	LogError(err, "getmsglist error")
	return msgs
}

func AddMsg(msg string) (bool, int) {
	id := setMaxId()
	fmt.Printf("get max id: %d\n", id)
	mid := strconv.Itoa(id)
	msg = mid + "@$@" + msg
	ok, err := client.Zadd(config.BoardCastKey, []byte(msg), float64(id))
	CheckError(err, "add msg error")
	return ok, id
}

func setMaxId() int {
	id := GetMsg("MSG_MAX_ID")
	mid := 0
	if nil != id {
		mid, _ = strconv.Atoi(string(id))
	}
	mid++
	err := client.Set("MSG_MAX_ID", []byte(strconv.Itoa(mid)))
	CheckError(err, "set max id error")
	return mid
}

func SetSendMaxId() {
	msgs, err := client.Zrevrange(config.BoardCastKey, 0, 0)
	LogError(err, "set send max id error")
	id := []byte("0")
	if len(msgs) > 0 {
		arrs := bytes.Split(msgs[0], []byte("@$@"))
		if len(arrs) > 1 {
			id = arrs[0]
		}
	}

	client.Set("SEND_MSG_MAX_ID", id)
}

func GetSendMaxId() int {
	id := GetMsg("SEND_MSG_MAX_ID")
	mid := 0
	if nil != id {
		mid, _ = strconv.Atoi(string(id))
	}
	return mid
}

func GetData(key string) []byte {
	val, _ := client.Get(key)
	return val
}

func SetData(msg string) {

}

func SaveData(key string, val []byte) bool {
	err := client.Set(key, val)
	if err != nil {
		return false
	}
	return true
}

func GetMd5(key string) string {
	h := md5.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}

func MergeByte(b1, b2 []byte) []byte {
	tmp := make([][]byte, 0)
	tmp[0] = b1
	tmp[1] = b2

	return bytes.Join(tmp, []byte{})
}
