package roomPackage

import "chatServer/protocol"

type RoomConfig struct {
	StartRoomNumber	int32
	MaxRoomCount	int32
	MaxUserCount	int32
}


type RoomUser struct {
	NetSessionIndex		int32
	NetSessionUniqueID	uint64

	RoomUniqueID		uint64
	IDLen				int8
	ID					[protocol.MAX_USER_ID_BYTE_LENGTH]byte

	packetDataSize		int16
}


func (user *RoomUser)Init (userID []byte, uniqueID uint64) {
	IDLen := len(userID)

	user.IDLen = int8(IDLen)
	copy(user.ID[:], userID)

	user.RoomUniqueID = uniqueID
}


func (user *RoomUser) SetNetworkInfo(sessionIndex int32, sessionUniqueID uint64) {
	user.NetSessionIndex = sessionIndex
	user.NetSessionUniqueID = sessionUniqueID
}

func (user *RoomUser) PacketDataSize() int16 {
	return int16(1) + int16(user.IDLen) + 8
}

type AddRoomUserInfo struct {
	userID []byte

	NetSessionIndex		int32
	NetSessionUniqueID	uint64
}