package main

import (
    "fmt"
    "miop"
    "bytes"
)


func main() {
    fmt.Printf("===========start 1==================\n")
    emptyProp()
    fmt.Printf("===========end 1==================\n")
    fmt.Printf("===========start 2==================\n")
    emptyData()
    fmt.Printf("===========end 2==================\n")
    msg := miop.MiopDataPack(uint32(123), "i am message") 
    fmt.Printf("msg len %d \n", len(msg))
    b := miop.MiopPack("k1:v1\nu:shenzhe@vedio\nk3:v3", msg)
    fmt.Printf("len:%d, 16:%x\n", len(b), b)
    //b = bytes.Join(b, []byte)
    d := make([][]byte, 2)
    d[0] = b;
    d[1] = b;
    b = bytes.Join(d, nil)
    fmt.Printf("len:%d, 16:%x\n", len(b), b)
    fmt.Printf("\n===pack end===\n")

    datas, more := miop.MiopUnpack(b[:80])
    fmt.Printf("datas len %d \n", len(datas))
    for _, data := range datas { 
    fmt.Printf("%d, %d, %s, %d, %x", data.Ve, data.Pl, data.Prop, data.Dl, data.Data)
    fmt.Printf("\n===unpack end===\n")

    prop := data.ParseProp()

    u, p := prop.GetUser()

        fmt.Printf("u: %s, p: %s\n", u, p) 

        msgContent := data.ParseData()
    
        for _, v := range msgContent {
            fmt.Printf("\n message: %d, %d, %s \n", v.Id, v.Bl, v.Msg)
        }
    }

    fmt.Printf("more: %d, % x", len(more), more)
    d[0] = more
    d[1] = b[80:]
    b = bytes.Join(d, nil)
    fmt.Printf("len:%d, 16:%x\n", len(b), b)
    fmt.Printf("\n===pack end===\n")

    datas, more = miop.MiopUnpack(b)
    fmt.Printf("datas len %d \n", len(datas))
    for _, data := range datas {
        fmt.Printf("%d, %d, %s, %d, %x", data.Ve, data.Pl, data.Prop, data.Dl, data.Data)
        fmt.Printf("\n===unpack end===\n")

        prop := data.ParseProp()

        u, p := prop.GetUser()

        fmt.Printf("u: %s, p: %s\n", u, p)

        msgContent := data.ParseData()

        for _, v := range msgContent {
            fmt.Printf("\n message: %d, %d, %s \n", v.Id, v.Bl, v.Msg)
        }
    }
}

func emptyProp() {
    msg := miop.MiopDataPack(uint32(123), "i am message") 
    fmt.Printf("msg len %d \n", len(msg))
    b := miop.MiopPack("", msg)
    fmt.Printf("len:%d, 16:%x\n", len(b), b)
    datas, _ := miop.MiopUnpack(b[:])
    fmt.Printf("datas len %d \n", len(datas))
    for _, data := range datas {
        fmt.Printf("ve:%d, pl:%d, prop:%s, dl:%d, data:%x \n", data.Ve, data.Pl, data.Prop, data.Dl, data.Data)
        prop := data.ParseProp()

        u, p := prop.GetUser()

        fmt.Printf("u: %s, p: %s\n", u, p)

        msgContent := data.ParseData()

        for _, v := range msgContent {

            fmt.Printf("\n message: %d, %d, %s \n", v.Id, v.Bl, v.Msg)

        }
    }
    
}

func emptyData() {
    b := miop.MiopPack("k1:v1\nu:shenzhe@vedio\nk3:v3", []byte(""))
    fmt.Printf("len:%d, 16:%x\n", len(b), b)
    datas, _ := miop.MiopUnpack(b[:])
    fmt.Printf("datas len %d \n", len(datas))
    for _, data := range datas {
        fmt.Printf("ve:%d, pl:%d, prop:%s, dl:%d, data:%x \n", data.Ve, data.Pl, data.Prop, data.Dl, data.Data)

        prop := data.ParseProp()

        u, p := prop.GetUser()

        fmt.Printf("u: %s, p: %s\n", u, p)

        msgContent := data.ParseData()

        for _, v := range msgContent {

            fmt.Printf("\n message: %d, %d, %s \n", v.Id, v.Bl, v.Msg)

        }
    }
    
}
