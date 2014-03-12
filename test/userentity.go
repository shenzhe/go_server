package main

import (
    "entity"
    "fmt"
    "time"
    "config"
)

func main() {
    uid := "shenzhe@vedio"
    user := entity.GetUser(uid)
    nowTime := time.Now().Unix()

    if nil == user {
         fmt.Printf("%s empty\n", uid)
         user = &entity.UserEntity{uid, "127.0.0.1", 100000, nowTime, nowTime, nowTime, config.Room}

         if user.Save() {
             fmt.Printf("user save success\n")
         } else {
             fmt.Printf("user save err\n")
         }
     } else {
          user.LoginTime = nowTime
          user.Save()
          fmt.Printf("%s %s %d %d %d %d %s\n", user.Uid, user.Ip, user.CityCode, user.CreateTime, user.LoginTime, user.LastTime, user.Room)
     }
}


