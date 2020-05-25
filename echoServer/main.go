package main

import (
	"flag"
	NetLib "gohipernetFake"
)
func main() {

	NetLib.NetLibInitLog()

	// Config 파일에서 정보 가져오기
	config := NetLib.NetworkConfig{}

	flag.BoolVar(&config.IsTcp4Addr, "c_IsTcp4Addr", true, "bool flag")
	flag.StringVar(&config.BindAddress, "c_BindAddress", "127.0.0.1:11021", "string flag")
	flag.IntVar(&config.MaxSessionCount, "c_MaxSessionCount", 0, "int flag")
	flag.IntVar(&config.MaxPacketSize, "c_MaxPacketSize", 0, "int_flag")

	flag.Parse()


	//서버 시작
	StartServer(config)

}
