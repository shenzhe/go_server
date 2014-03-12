package main

import (
    "fmt"
    "comm"
)

func main() {
    u, k := "shenzhe@vedio", "adfip813"
    md5Str := comm.GetMd5(u+k)
 
    fmt.Printf("u:%s k:%s md5:%s", u, k, md5Str)

    if "685c59a3a1451ca2e68fa2262ae890bb" == md5Str {
        fmt.Printf("md5 check succ")
    }
}
