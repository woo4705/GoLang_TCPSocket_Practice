package roomPackage

import (
	"chatServer/protocol"
	"go.uber.org/zap"
	NetLib "gohipernetFake"
)

type RoomManager struct {
	RoomStartNum	int32
	MaxRoomCount	int32
	RoomCountList	[]int16
	RoomList		[]BaseRoom
}

func NewRoomManager(config RoomConfig) *RoomManager {
	roomManager := new(RoomManager)
	roomManager.Initialize(config)
	return roomManager
}


func (roomMgr *RoomManager) Initialize(config RoomConfig) {
	roomMgr.RoomStartNum = config.StartRoomNumber
	roomMgr.MaxRoomCount = config.MaxRoomCount
	roomMgr.RoomCountList = make([]int16, config.MaxRoomCount)
	roomMgr.RoomList = make([]BaseRoom, config.MaxRoomCount)

	for i:= int32(0); i < roomMgr.MaxRoomCount; i++ {
		roomMgr.RoomList[i].Initialize(i, config)
		roomMgr.RoomList[i].SettingPacketFunction()
	}

	LogStartRoomPacketProcess(config.MaxRoomCount, config)

	NetLib.NTELIB_LOG_INFO("[RoomManager Initialize]", zap.Int32("MaxRoomCount",roomMgr.MaxRoomCount))

}


func (roomMgr *RoomManager) GetAllChannelUserCount() []int16 {
	return nil
}


func (roomMgr *RoomManager) GetRoomByNumber(roomNumber int32) *BaseRoom {
	roomIndex := roomNumber - roomMgr.RoomStartNum

	if roomIndex < 0 || roomIndex >= roomMgr.MaxRoomCount {
		return nil
	}

	return &roomMgr.RoomList[roomIndex]
}


func (roomMgr *RoomManager) GetRoomUserCount(roomID int32) int32 {
	return roomMgr.RoomList[roomID].GetCurrentUserCount()
}


func (roomMgr *RoomManager) PacketProcess(roomNumber int32, packet protocol.Packet){
	NetLib.NTELIB_LOG_DEBUG("[RoomManager - PacketProcess]", zap.Int16("PacketID",packet.ID))
	isRoomEnterReq := false

	if roomNumber == -1 && packet.ID == protocol.PACKET_ID_ROOM_ENTER_REQ {
		isRoomEnterReq = true

		var requestPacket protocol.RoomEnterRequestPacket
		(&requestPacket).DecodingPacket(packet.Data)

		roomNumber = requestPacket.RoomNumber
	}

	room := roomMgr.GetRoomByNumber(roomNumber)
	if room == nil {
		protocol.NotifyErrorPacket(packet.UserSessionIndex, packet.UserSessionUniqueID, protocol.ERROR_CODE_USER_NOT_IN_ROOM)
		return
	}

	user := room.GetUser(packet.UserSessionUniqueID)
	if user == nil && isRoomEnterReq == false {
		protocol.NotifyErrorPacket(packet.UserSessionIndex, packet.UserSessionUniqueID, protocol.ERROR_CODE_USER_NOT_IN_ROOM)
		return
	}

	funcCount := len(room.FuncPacketIDList)

	for i:=0; i<funcCount; i++ {
		if room.FuncPacketIDList[i] != packet.ID {
			continue
		}

		result := room.FuncList[i] (user, packet)
		if result != protocol.ERROR_CODE_NONE {
			NetLib.NTELIB_LOG_DEBUG("[Room - PacketProcess Fail]",
			zap.Int16("PacketID",packet.ID), zap.Int16("Error", result))
		}

		return
	}

	NetLib.NTELIB_LOG_DEBUG("[Room- PacketProcess - Fail(Not Registered)]", zap.Int16("PacketID", packet.ID))
}



func LogStartRoomPacketProcess(maxRoomCount int32, config RoomConfig) {
	NetLib.NTELIB_LOG_INFO("[RoomManager startRoomPacketProcess]",
		zap.Int32("maxRoomCount", maxRoomCount),
		zap.Int32("StartRoomNumber", config.StartRoomNumber),
		zap.Int32("MaxUserCount", config.MaxUserCount),
	)
}