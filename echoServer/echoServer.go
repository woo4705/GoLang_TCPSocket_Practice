package main

import (
	NetLib "gohipernetFake"
	"strconv"
	"strings"
	"go.uber.org/zap"
)

type echoServer struct {
	IP		string
	Port	int
}



func StartServer(netConfig NetLib.NetworkConfig) {
	server := echoServer{}

	if( server.setIPAddr (netConfig.BindAddress) ==false){
		NetLib.NTELIB_LOG_ERROR("Server IPAddr Bind Error")
	}


	//TODO:네트워크 라이브러리의 함수들 맵핑해주기. 함수들 본체 구현해주기
	sessionNetFunctors := NetLib.SessionNetworkFunctors{}
	sessionNetFunctors.OnConnect = server.OnConnect
	sessionNetFunctors.OnReceive = server.OnReceive
	sessionNetFunctors.OnClose = server.OnClose
	sessionNetFunctors.OnReceiveBufferedData = nil
	sessionNetFunctors.PacketTotalSizeFunc = NetLib.PacketTotalSize
	sessionNetFunctors.PacketHeaderSize = NetLib.PACKET_HEADER_SIZE
	sessionNetFunctors.IsClientSession = true

	NetLib.NetLibInitNetwork(NetLib.PACKET_HEADER_SIZE, NetLib.PACKET_HEADER_SIZE)
	NetLib.NetLibStartNetwork(&netConfig,sessionNetFunctors)


}



func (server *echoServer)setIPAddr(ipAddr string) bool{
	result := strings.Split(ipAddr,":")
	if len(result) != 2 {
		return false
	}

	server.IP = result[0]
	server.Port,_ = strconv.Atoi(result[1])

	NetLib.NTELIB_LOG_INFO("IP Addr:", zap.String("IP", server.IP), zap.Int("Port", server.Port) )
	return true

}





func (server *echoServer)OnConnect(sessionIdx int32, sessionID uint64) {
	NetLib.NTELIB_LOG_INFO("Server Connected", zap.Int32("sessionIdx",sessionIdx), zap.Uint64("sessionID",sessionID) )
}


func (server *echoServer)OnReceive(sessionIdx int32, sessionID uint64, data[] byte) bool {
	//데이터를 전송
	NetLib.NTELIB_LOG_INFO("Server Connected", zap.Int32("sessionIdx",sessionIdx), zap.Uint64("sessionID",sessionID) )
	NetLib.NetLibISendToClient(sessionIdx, sessionID, data)
	return true

}


func (server *echoServer)OnClose(sessionIdx int32, sessionID uint64) {
	NetLib.NTELIB_LOG_INFO("Server Disconnected", zap.Int32("sessionIdx",sessionIdx), zap.Uint64("sessionID",sessionID) )

}

