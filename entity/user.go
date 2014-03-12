package entity

import (
    "comm"
    "encoding/json"
    "config"
)

type UserEntity struct {
    Uid        string
    Ip         string
    CityCode   int
    CreateTime int64
    LoginTime  int64
    LastTime   int64
    Room       string
}

func GetEntity() *UserEntity{
    return &UserEntity{"", "", 0, 0, 0, 0, config.Room}
}

type UserFilter struct {
    CityCode int
    LastTime int64
}

func (ue *UserEntity) Save() bool {
    data, err := json.Marshal(ue)
    comm.CheckError(err, "ue json encode error")
    return comm.SaveData(ue.Uid, data)
}

func GetUser(uid string) *UserEntity {
    data := comm.GetData(uid)
    if nil == data {
        return nil
    }

    var user UserEntity
    err := json.Unmarshal(data, &user)

    if nil != err {
        return nil
    }

    return &user
}

func AddUser(uid, ip string, cc int, ct, lt, lat int64) *UserEntity {
    user := &UserEntity{uid, ip, cc, ct, lt, lat, config.Room}
    if user.Save() {
        return user
    }

    return nil
}

func (uf *UserFilter) CheckFilter(ue UserEntity) bool {
    if uf.CityCode > 0 {
        if uf.CityCode != ue.CityCode {
            return false
        }
    }
    if uf.LastTime > 0 {
        if uf.LastTime < ue.LastTime {
            return false
        }
    }

    return true

}
