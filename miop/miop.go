package miop

import (
    "bytes"
    "comm"
    "config"
    "encoding/binary"
    "encoding/json"
    "entity"
    "strconv"
    "strings"
)

type MiopData struct {
    Id  uint32
    Bl  uint32
    Msg string
}

type MiopProp struct {
    Prop map[string]string
}
type Miop struct {
    Ve   uint16
    Pl   uint16
    Dl   uint32
    Prop []byte
    Data []byte
}

func MiopDataPack(id uint32, content string) []byte {
    buf := new(bytes.Buffer)
    var data = []interface{}{
        id,
        uint32(len(content)),
        []byte(content),
    }
    for i, v := range data {
        err := binary.Write(buf, binary.BigEndian, v)
        if err != nil {
            comm.CheckError(err, "miopdata binary.Write failed:"+string(i))
        }
    }

    return buf.Bytes()
}

func MiopPack(prop string, msg []byte) []byte {
    buf := new(bytes.Buffer)
    var data = []interface{}{
        uint16(0),
        uint16(len(prop)),
        uint32(len(msg)),
        []byte(prop),
        msg,
    }

    for i, v := range data {
        err := binary.Write(buf, binary.BigEndian, v)
        if err != nil {
            comm.CheckError(err, "binary.Write failed:"+string(i))
        }
    }

    return buf.Bytes()
}

func MiopUnpack(b []byte) ([]*Miop, []byte) {
    blen := len(b)
    mlen := blen
    comm.Debug("blen: %d, mlen: %d\n", blen, mlen)
    buff := bytes.NewReader(b)
    miops := make([]*Miop, 0)
    more := make([]byte, 0)
    for {
        if blen <= 0 {
            break
        }
        var ve uint16
        err := binary.Read(buff, binary.BigEndian, &ve)
        comm.CheckError(err, "read ve error")
        var pl uint16
        err = binary.Read(buff, binary.BigEndian, &pl)
        comm.CheckError(err, "read pl error")
        var dl uint32
        err = binary.Read(buff, binary.BigEndian, &dl)
        comm.CheckError(err, "read dl error")
        msgLen := int(8 + uint32(pl) + dl)
        comm.Debug("pl:%d, dl:%d, msgLen: %d, blen: %d\n", pl, dl, msgLen, blen)
        if blen < msgLen {
            break
        }
        prop := make([]byte, pl)
        err = binary.Read(buff, binary.BigEndian, &prop)
        comm.CheckError(err, "read prop error")
        data := make([]byte, dl)
        err = binary.Read(buff, binary.BigEndian, &data)
        comm.CheckError(err, "read data error")

        miops = append(miops, &Miop{ve, pl, dl, prop, data})

        blen -= msgLen

    }

    if blen > 0 {
        ml := mlen - blen
        more = b[ml:]
    }

    return miops, more
}

func (miop *Miop) ParseProp() *MiopProp {
    if miop.Pl < 1 {
        return nil
    }
    arrs := bytes.Split(miop.Prop, []byte("\n"))
    propMap := make(map[string]string)
    for _, v := range arrs {
        arr := bytes.Split(v, []byte(":"))
        if len(arr) > 1 {
            propMap[string(arr[0])] = string(arr[1])
        }
    }

    return &MiopProp{propMap}
}

func (prop *MiopProp) GetPropKey(key string) string {
    if nil == prop {
        return ""
    }
    val, ok := prop.Prop[key]
    if ok {
        return val
    }
    return ""
}

func (prop *MiopProp) GetLastId() (int) {
    val :=  prop.GetPropKey("lid")
    if "" == val {
        return -1
    }
    
    id, _ := strconv.Atoi(val)
    return id
}
func (prop *MiopProp) GetUser() (string, string) {
    val := prop.GetPropKey("u")
    if "" == val {
        return "", ""
    }

    infos := strings.Split(val, "@")

    u := infos[0]

    p := ""

    if len(infos) > 1 {
        p = infos[1]
    }
    return u, p
}

func (miop *Miop) ParseData() []*MiopData {
    if miop.Dl < 1 {
        return nil
    }
    buff := bytes.NewReader(miop.Data)
    arr := make([]*MiopData, 0)
    dlen := len(miop.Data)
    for {
        if dlen <= 0 {
            break
        }
        var id uint32
        err := binary.Read(buff, binary.BigEndian, &id)
        comm.CheckError(err, "read id error")
        if id <= 0 {
            break
        }

        var dl uint32
        err = binary.Read(buff, binary.BigEndian, &dl)
        comm.CheckError(err, "read dl error")
        var content = make([]byte, dl)
        err = binary.Read(buff, binary.BigEndian, &content)
        comm.CheckError(err, "read content error")

        arr = append(arr, &MiopData{id, dl, string(content)})

        dlen -= int(8 + dl)
    }

    return arr
}

func (prop *MiopProp) CheckToken() bool {
    token := prop.GetPropKey("token")
    if "" == token {
        comm.Log("token empty")
        return false
    }

    if token == config.AdminToken {
        return true
    }

    u := prop.GetPropKey("u")
    if "" == u {
        comm.Log("u empty")
        return false
    }

    checkToken := comm.GetMd5(u + config.Key)

    comm.Log("check token:%s checkToken: %s", token, checkToken)
    return token == checkToken

}

func (prop *MiopProp) GetFilter() *entity.UserFilter {
    filter := prop.GetPropKey("filter")

    if "" == filter {
        return nil
    }

    var userfilter entity.UserFilter

    err := json.Unmarshal([]byte(filter), &userfilter)

    if nil != err {
        return nil
    }

    return &userfilter
}

func ParseMsgList(msgs [][]byte) [][]byte {
    msgpack := make([][]byte, 0)
    for _, msg := range msgs {
        arrs := bytes.Split(msg, []byte("@$@")) 
        if len(arrs) > 1 {
            id, _ := strconv.Atoi(string(arrs[0]))
            msgpack = append(msgpack, MiopDataPack(uint32(id), string(arrs[1])))
        }
    }
    return msgpack
}
