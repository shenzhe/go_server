package config

import (
	"flag"
)

var (
	Version             = 0
	HostAndPort         = ":9090"
	BuffLen             = 1024
	HeartbeatTime int64 = 10
	LogPath             = "/data/go/log/"
	Key                 = "adfop013@dafa"
	Deamonize           = false
	Room                = "127.0.0.1:9090"
	BoardCastKey        = "boardcast_movie"

	RoomList = make(map[string]int)
)

var (
	RLFile         = "/data/go/src/conf/dispatcher.conf"
	AdminToken     = "piuer813#@80"
	AdminUid       = "admin"
	DispatcherHost = "127.0.0.1:9099"
)

func Init() {
	flag.StringVar(&HostAndPort, "hp", "0.0.0.0:9090", "host and port")
	flag.StringVar(&LogPath, "log", "/data/go/log/", "lop path")
	flag.BoolVar(&Deamonize, "d", false, "is deamonize")
	flag.Parse()
}
