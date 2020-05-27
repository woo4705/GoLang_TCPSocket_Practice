package roomPackage

import (
	"chatServer/connectedSession"
	"chatServer/protocol"
	"go.uber.org/zap"
	NetLib "gohipernetFake"

)

func (room *BaseRoom) PacketProcess_Relay(user *RoomUser, packet protocol.Packet) int16 {
	var relayNotify protocol.RoomRelayNotifyPacket
	relayNotify.RoomUserUniqueID = user.RoomUniqueID
	relayNotify.Data = packet.Data

	notifySendBuf, packetSize := relayNotify.EncodingPacket(packet.DataSize)
	room.BroadCastPacket(packetSize, notifySendBuf, 0)

	NetLib.NTELIB_LOG_DEBUG("Room Relay", zap.String("Sender",string(user.ID[:])) )
	return protocol.ERROR_CODE_NONE

}



func (room *BaseRoom) PacketProcess_EnterUser(inValidUser *RoomUser, packet protocol.Packet) int16 {
	curTime := NetLib.NetLib_GetCurrnetUnixTime()
	sessionIndex := packet.UserSessionIndex
	sessionUniqueID := packet.UserSessionUniqueID

	NetLib.NTELIB_LOG_INFO("[Room PacketProcess EnterUser]")

	var requestPacket protocol.RoomEnterRequestPacket
	(&requestPacket).DecodingPacket(packet.Data)

	userID, ok := connectedSession.GetUserID(sessionIndex)
	if ok == false {
		SendRoomEnterResult(sessionIndex, sessionUniqueID, 0,0,protocol.ERROR_CODE_ENTER_ROOM_INVALID_USER_ID)
		return protocol.ERROR_CODE_ENTER_ROOM_INVALID_USER_ID
	}

	userInfo := AddRoomUserInfo{
		userID,
		sessionIndex,
		sessionUniqueID,
	}
	newUser, addResult := room.AddUser(userInfo)


	if addResult != protocol.ERROR_CODE_NONE {
		SendRoomEnterResult(sessionIndex, sessionUniqueID, 0,0, addResult)
		return addResult
	}

	if connectedSession.SetRoomNumber(sessionIndex, sessionUniqueID, room.GetNumber(), curTime) == false {
		SendRoomEnterResult(sessionIndex, sessionUniqueID, 0, 0, protocol.ERROR_CODE_ENTER_ROOM_INVALID_SESSION_STATE)
		return protocol.ERROR_CODE_ENTER_ROOM_INVALID_SESSION_STATE
	}


	if room.GetCurrentUserCount() > 1 {
		room.SendNewUserInfoPacket(newUser)
		room.SendUserInfoListPacket(newUser)
	}

	roomNumber := room.GetNumber()
	SendRoomEnterResult(sessionIndex, sessionUniqueID, roomNumber, newUser.RoomUniqueID, protocol.ERROR_CODE_NONE)

	return protocol.ERROR_CODE_NONE

}



func SendRoomEnterResult(sessionIndex int32, sessionUniqueID uint64, roomNumber int32, userUniqueID uint64, result int16){
	response := protocol.RoomEnterResponsePacket{
		result,
		roomNumber,
		userUniqueID,
	}

	sendPacket,_ := response.EncodingPacket()
	NetLib.NetLibIPostSendToClient(sessionIndex, sessionUniqueID, sendPacket)
}



func (room *BaseRoom) SendUserInfoListPacket(user *RoomUser){
	NetLib.NTELIB_LOG_DEBUG("Room SendUserInfoListPacket", zap.Uint64("SessionUniqueID", user.NetSessionUniqueID))

	userCount, userInfoListSize, userInfoListBuffer := room.AllocAllUserInfo(user.NetSessionUniqueID)

	var response protocol.RoomUserListNotifyPacket
	response.UserCount = userCount
	response.UserList = userInfoListBuffer
	sendBuf, _ := response.EncodingPacket(userInfoListSize)

	NetLib.NetLibIPostSendToClient(user.NetSessionIndex, user.NetSessionUniqueID, sendBuf)
}


func (room *BaseRoom) SendNewUserInfoPacket(user *RoomUser){
	NetLib.NTELIB_LOG_DEBUG("Room SendNewUserInfoPakcet", zap.Uint64("SessionUniqueID", user.NetSessionUniqueID))

	userInfoSize, userInfoListBuffer := room.AllocUserInfo(user)

	var response protocol.RoomNewUserNotifyPacket
	response.User = userInfoListBuffer
	sendBuf, packetSize := response.EncodingPacket(userInfoSize)

	room.BroadCastPacket(int16(packetSize), sendBuf, user.NetSessionUniqueID)

}




func (room *BaseRoom) PacketProcess_LeaveUser(user *RoomUser, packet protocol.Packet) int16{
	NetLib.NTELIB_LOG_DEBUG("[Room packetProcess_LeaveUser]")

	room.LeaveUserProcess(user)

	userSessionIndex := user.NetSessionIndex
	userSessionUniqueID := user.NetSessionUniqueID

	SendRoomLeaveResult(userSessionIndex, userSessionUniqueID, protocol.ERROR_CODE_NONE )

	return protocol.ERROR_CODE_NONE
}


func (room *BaseRoom) LeaveUserProcess(user *RoomUser)  {
	NetLib.NTELIB_LOG_DEBUG("[Room LeaveUser Processs]")

	roomUserUniqueID := user.RoomUniqueID
	userSessionIndex := user.NetSessionIndex
	userSessionUniqueID := user.NetSessionUniqueID

	room.RemoveUser(user)
	room.SendRoomLeaveUserNotify(roomUserUniqueID,userSessionUniqueID )

	curTime := NetLib.NetLib_GetCurrnetUnixTime()
	connectedSession.SetRoomNumber(userSessionIndex, userSessionUniqueID, -1, curTime)

}


func SendRoomLeaveResult(sessionIndex int32, sessionUniqueID uint64, result int16){
	response := protocol.RoomLeaveUserResponsePacket{Result: result}
	sendPacket,_ := response.EncodingPacket()
	NetLib.NetLibIPostSendToClient(sessionIndex, sessionUniqueID, sendPacket)
}


func (room *BaseRoom) SendRoomLeaveUserNotify(roomUserUniqueID uint64, userSessionUniqueID uint64) {
	NetLib.NTELIB_LOG_DEBUG("Room SendRoomLeaveUserNotify", zap.Uint64("SessionUniqueID", userSessionUniqueID), zap.Int32("RoomIndex", room.Index) )

	notifyPacket := protocol.RoomLeaveUserNotifyPacket{UserUniqueID: roomUserUniqueID}
	sendBuf, packetSize := notifyPacket.EncodingPacket()
	room.BroadCastPacket(int16(packetSize), sendBuf, userSessionUniqueID)
}





func (room *BaseRoom) PacketProcess_Chat(user *RoomUser, packet protocol.Packet) int16{
	sessionIndex := packet.UserSessionIndex
	sessionUniqueID := packet.UserSessionUniqueID

	var chatPacket protocol.RoomChatRequestPacket
	if chatPacket.Decoding(packet.Data) == false {
		SendRoomChatResult(sessionIndex, sessionUniqueID, protocol.ERROR_CODE_PACKET_DECODING_FAIL)
		return protocol.ERROR_CODE_PACKET_DECODING_FAIL
	}


	msgLen := len(chatPacket.MsgData)
	if msgLen < 1 || msgLen > protocol.MAX_CHAT_MESSAGE_BYTE_LENGTH {
		SendRoomChatResult(sessionIndex, sessionUniqueID, protocol.ERROR_CODE_ROOM_CHAT_CHAT_MSG_LEN)
		return protocol.ERROR_CODE_ROOM_CHAT_CHAT_MSG_LEN
	}

	var chatNotifyResponse protocol.RoomChatNotifyPacket
	chatNotifyResponse.RoomUserUniqueID = user.RoomUniqueID
	chatNotifyResponse.MsgLen = int16(msgLen)
	chatNotifyResponse.Msg = chatPacket.MsgData

	notifySendBuf, packetSize := chatNotifyResponse.EncodingPacket()
	room.BroadCastPacket(packetSize, notifySendBuf, 0)

	SendRoomChatResult(sessionIndex, sessionUniqueID, protocol.ERROR_CODE_NONE)

	NetLib.NTELIB_LOG_DEBUG("Channel Chat Notify Function", zap.String("Sender", string(user.ID[:])),
		zap.String("Message", string(chatPacket.MsgData) ))


	return protocol.ERROR_CODE_NONE
}



func SendRoomChatResult(sessionIndex int32, sessionUniqueID uint64, result int16) {
	response := protocol.RoomChatRequestPacket{MsgLength: result}
	sendPacket,_ := response.EncodingPacket()
	NetLib.NetLibIPostSendToClient(sessionIndex, sessionUniqueID, sendPacket)
}