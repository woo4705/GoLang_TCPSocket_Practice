package main
import (
	"flag"
	NetLib "gohipernetFake"
)

func main() {
	NetLib.NetLibInitLog()

	// Config 파일에서 정보 가져오기
	netConfig, appConfig := ParseConfigData()
	netConfig.WriteNetworkConfig(true)

	//서버 시작
	CreateAndStartServer(netConfig, appConfig)
}

func ParseConfigData() (NetLib.NetworkConfig, ConfigAppServer)  {
	NetLib.NTELIB_LOG_INFO("[Setting NetworkConfig]")

	appConfig := ConfigAppServer{
		"chatServer",
		1000,
		0,
		4,
	}

	netConfig := NetLib.NetworkConfig{}

	flag.BoolVar(&netConfig .IsTcp4Addr, "c_IsTcp4Addr", true, "bool flag")
	flag.StringVar(&netConfig .BindAddress, "c_BindAddress", "127.0.0.1:11021", "string flag")
	flag.IntVar(&netConfig .MaxSessionCount, "c_MaxSessionCount", 0, "int flag")
	flag.IntVar(&netConfig .MaxPacketSize, "c_MaxPacketSize", 0, "int_flag")

	flag.Parse()

	return  netConfig, appConfig

}