package main

import (
    "math/rand"
    "fmt"
    "time"
)

func main() {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i:=0; i<10; i++ {
        fmt.Println(r.Intn(2))
    }

}
