package main

import (
	"chatServer/protocol"
	"go.uber.org/zap"
	NetLib "gohipernetFake"
)

func (server *ChatServer) DistributePacket(sessionIndex int32, sessionUniqueID uint64, packetData []byte){
	packetID := protocol.PeekPacketID(packetData)
	bodySize, bodyData := protocol.PeekPacketBody(packetData)

	NetLib .NTELIB_LOG_DEBUG("Distribute Packet",
		zap.Int32("sessionIndex", sessionIndex), zap.Uint64("sessionUniqueID",sessionUniqueID), zap.Int16("PacketID",packetID))

	packet := protocol.Packet{ID: packetID}
	packet.UserSessionIndex = sessionIndex
	packet.UserSessionUniqueId = sessionUniqueID
	packet.ID = packetID
	packet.DataSize = bodySize
	packet.Data = make([]byte, packet.DataSize)
	copy(packet.Data, bodyData)



}


func (server *ChatServer) PacketProcess_goroutine() {

}


func (server *ChatServer) PacketProcess_goroutine_Impl() bool {
	IsWantTermination := false

	return IsWantTermination
}

func ProcessPacketLogin(sessionIndex int32, bodySize int16, bodyData []byte)  {

}

func SendLoginResult(sessionIndex int32, sessionUniqueID uint64, result int16){

}

func ProcessPacketSesssionClosed(server *ChatServer, sessionIndex int32, sessionUnique uint64){

}