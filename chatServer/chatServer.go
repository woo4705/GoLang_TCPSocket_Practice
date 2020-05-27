package main

import (
	"chatServer/connectedSession"
	"chatServer/protocol"
	"chatServer/roomPackage"
	"go.uber.org/zap"
	NetLib "gohipernetFake"
	"strconv"
	"strings"
)

type ConfigAppServer struct {
	GameName			string

	RoomMaxCount		int32
	RoomStartNum		int32
	RoomMaxUserCount	int32
}

type ChatServer struct {
	ServerIndex		int
	IP				string
	Port			int

	PacketChannel	chan protocol.Packet
	RoomMgr 		roomPackage.RoomManager
}


func CreateAndStartServer(netConfig NetLib.NetworkConfig, appConfig ConfigAppServer){
	NetLib.NTELIB_LOG_INFO("Create Server")
	var server ChatServer

	if server.setIPAddress(netConfig.BindAddress) == false {
		NetLib.NTELIB_LOG_ERROR("Server create fail. IP Address")
		return
	}

	protocol.Init_packetSize()

	maxUserCount := appConfig.RoomMaxCount * appConfig.RoomMaxUserCount
	connectedSession.Init(int32(netConfig.MaxSessionCount), maxUserCount)

	server.PacketChannel = make(chan protocol.Packet, 256)

	/*
	Room packet부분 생성하기
	 */

	go server.PacketProcess_goroutine()


	networkFunctor := NetLib.SessionNetworkFunctors{}
	networkFunctor.OnConnect = server.OnConnect
	networkFunctor.OnReceive = server.OnReceive
	networkFunctor.OnClose = server.OnClose
	networkFunctor.OnReceiveBufferedData = nil
	networkFunctor.PacketTotalSizeFunc = NetLib.PacketTotalSize
	networkFunctor.PacketHeaderSize = NetLib.PACKET_HEADER_SIZE
	networkFunctor.IsClientSession = true

	NetLib.NetLibInitNetwork(NetLib.PACKET_HEADER_SIZE, NetLib.PACKET_HEADER_SIZE)
	NetLib.NetLibStartNetwork(&netConfig, networkFunctor)
}


func (server *ChatServer) setIPAddress(ipAddr string) bool {
	result := strings.Split(ipAddr,":")
	if len(result) != 2 {
		return false
	}

	server.IP = result[0]
	server.Port,_ = strconv.Atoi(result[1])

	NetLib.NTELIB_LOG_INFO("IP Addr:", zap.String("IP", server.IP), zap.Int("Port", server.Port) )
	return true
}


func (server *ChatServer) OnConnect(sessionIndex int32, sessionUniqueID uint64){
	NetLib.NTELIB_LOG_INFO("client OnConnect",
		zap.Int32("sessionIndex",sessionIndex), zap.Uint64("sessionUniqueID", sessionUniqueID),
		)
	connectedSession.AddSession(sessionIndex, sessionUniqueID )
}

func (server *ChatServer) OnReceive( sessionIndex int32, sessionUniqueID uint64, data[] byte ) bool {
	NetLib.NTELIB_LOG_INFO("OnReceive",
		zap.Int32("sessionIndex",sessionIndex), zap.Uint64("sessionUniqueID", sessionUniqueID), zap.Int("packetSize",len(data)),
		)
	server.DistributePacket(sessionIndex, sessionUniqueID, data)

	return true
}

func (server *ChatServer) OnClose( sessionIndex int32, sessionUniqueID uint64 ) {
	NetLib.NTELIB_LOG_INFO("client OnClose",
		zap.Int32("sessionIndex",sessionIndex), zap.Uint64("sessionUniqueID", sessionUniqueID),
	)
	server.DisconnectClient(sessionIndex, sessionUniqueID)
}


func (server* ChatServer) DisconnectClient( sessionIndex int32, sessionUniqueID uint64 )  {
	if connectedSession.IsLoginUser(sessionIndex) == false {
		NetLib.NTELIB_LOG_INFO("DisConnectClient", zap.Int32("sessionIndex",sessionIndex))
		connectedSession.RemoveSession(sessionIndex, false)
		return
	}

	packet := protocol.Packet{
		sessionIndex,
		sessionUniqueID,
		protocol.PACKET_ID_SESSION_CLOSE_SYS,
		0,
		nil,
	}

	server.PacketChannel <- packet

	NetLib.NTELIB_LOG_INFO("DisConnectClient Login User", zap.Int32("sessionIndex",sessionIndex))
}

