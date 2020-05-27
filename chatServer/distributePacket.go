package main

import (
	"bytes"
	"chatServer/connectedSession"
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
	packet.UserSessionUniqueID = sessionUniqueID
	packet.ID = packetID
	packet.DataSize = bodySize
	packet.Data = make([]byte, packet.DataSize)
	copy(packet.Data, bodyData)

	server.PacketChannel <- packet

	NetLib.NTELIB_LOG_DEBUG("distributePacket", zap.Int32("sessionIndex",sessionIndex), zap.Int16("PacketID",packetID))
}


func (server *ChatServer) PacketProcess_goroutine() {
	NetLib.NTELIB_LOG_DEBUG("start PacketProcess goroutine")

	for{
		if server.PacketProcess_goroutine_Impl() {
			NetLib.NTELIB_LOG_INFO("Wanted Stop PacketProcess goroutine")
			break
		}
	}

	NetLib.NTELIB_LOG_INFO("Stop rooms PacketProcess goroutine")
}


func (server *ChatServer) PacketProcess_goroutine_Impl() bool {
	IsWantTermination := false
	defer NetLib.PrintPanicStack()

	for{
		packet := <- server.PacketChannel
		sessionIndex := packet.UserSessionIndex
		sessionUniqieID := packet.UserSessionUniqueID
		bodySize := packet.DataSize
		bodyData := packet.Data

		if packet.ID == protocol.PACKET_ID_LOGIN_REQ {
			ProcessPacketLogin(sessionIndex, sessionUniqieID, bodySize, bodyData)
		}else if packet.ID == protocol.PACKET_ID_SESSION_CLOSE_SYS {
			ProcessPacketSesssionClosed(server, sessionIndex, sessionUniqieID)
		}else {
			roomNumber,_ := connectedSession.GetRoomNumber(sessionIndex)
			server.RoomMgr.PacketProcess(roomNumber, packet)
		}
	}

	return IsWantTermination
}



func ProcessPacketLogin(sessionIndex int32, sessionUniqueID uint64, bodySize int16, bodyData []byte)  {
	var reqPacket protocol.LoginRequestPacket

	if (&reqPacket).DecodingPacket(bodyData) == false {
		SendLoginResult(sessionIndex, sessionUniqueID, protocol.ERROR_CODE_PACKET_DECODING_FAIL)
		return
	}

	userID := bytes.Trim(reqPacket.UserID[:], "\x00")

	if len(userID) <= 0 {
		SendLoginResult(sessionIndex, sessionUniqueID, protocol.ERROR_CODE_LOGIN_USER_INVALID_ID)
		return
	}

	curTime := NetLib.NetLib_GetCurrnetUnixTime()

	if connectedSession.SetLogin(sessionIndex, sessionUniqueID, userID, curTime) == false {
		SendLoginResult(sessionIndex, sessionUniqueID, protocol.ERROR_CODE_LOGIN_USER_DUPLICATION)
		return
	}

	SendLoginResult(sessionIndex, sessionUniqueID, protocol.ERROR_CODE_NONE)
}



func SendLoginResult(sessionIndex int32, sessionUniqueID uint64, result int16){
	var resPacket protocol.LoginResponsePacket
	resPacket.Result = result
	sendPacket,_ := resPacket.EncodingPacket()

	NetLib.NetLibIPostSendToClient(sessionIndex, sessionUniqueID, sendPacket)
	NetLib.NTELIB_LOG_DEBUG("SendLoginResult", zap.Int32("sessionIndex", sessionIndex), zap.Int16("result",result))

}



func ProcessPacketSesssionClosed(server *ChatServer, sessionIndex int32, sessionUniqueID uint64){
	roomNumber,_ := connectedSession.GetRoomNumber(sessionIndex)

	if roomNumber > -1 {
		packet := protocol.Packet{
			sessionIndex,
			sessionUniqueID,
			protocol.PACKET_ID_ROOM_LEAVE_REQ,
			0,
			nil,
		}

		server.RoomMgr.PacketProcess(roomNumber, packet)
	}

	connectedSession.RemoveSession(sessionIndex, true)
}