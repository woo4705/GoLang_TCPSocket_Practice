package roomPackage

import (
	NetLib "gohipernetFake"
	"chatServer/protocol"
	"sync"
	"sync/atomic"
)

type BaseRoom struct {
	Index					int32
	Number					int32
	Config					RoomConfig

	CurUserCount			int32
	RoomUserUniqueIdSeq		uint64
	UserPool				*sync.Pool

	UserSessionUniqueIDMap	map[uint64]*RoomUser
	FuncPacketIDList		[]int16
	FuncList				[]func(*roomUser, protocol.Packet) int16

	EnterUserNotify			func(int64, int32)
	LeaveUserNotify			func(int64)
}


func (room *BaseRoom) GetIndex() int32 {
	return room.Index
}

func (room *BaseRoom) GetNumber() int32 {
	return room.Number
}

func (room *BaseRoom) GetCurrentUserCount() int32   {
	count := atomic.LoadInt32(&room.CurUserCount)
	return count
}

func (room *BaseRoom) GenerateUserUniqueID() uint64 {
	room.RoomUserUniqueIdSeq++
	uniqueID := room.RoomUserUniqueIdSeq
	return uniqueID
}



func (room *BaseRoom) Initialize(index int32, config RoomConfig){
	room.Number = config.StartRoomNumber
	room.Index = index
	room.Config = config

	room.InitUserPool()
	room.UserSessionUniqueIDMap = make(map[uint64]*RoomUser)
}



func (room *BaseRoom) EnableEnterUser() bool {
	if room.GetCurrentUserCount() >= room.Config.MaxRoomCount {
		return false
	}
	return true
}


func (room *BaseRoom) SettingPacketFunction() {
	
}

func (room *BaseRoom) AddPacketFunction(packetID int16, packetFunc func(*RoomUser, protocol.Packet) int16){

}



func (room *BaseRoom) InitUserPool() {

}

func (room *BaseRoom) GetUserObject() *RoomUser {
	return nil
}

func (room *BaseRoom) PutUserObject(user *RoomUser) {

}


func (room *BaseRoom) IsFullUser() bool {
	return false
}

func (room *BaseRoom) RemoveUser(user *RoomUser) {

}

func (room *BaseRoom) RemoveUserObject(user *RoomUser) {

}

func (room *BaseRoom) GetUser(sessionUniqueID uint64) *RoomUser {
	return nil
}

func (room *BaseRoom) AllocAllUserInfo(exceptSessionUniqueID uint64) (userCount int8, dataSize int16, dataBuffer []byte) {
	return userCount, dataSize, dataBuffer
}



func WriteUserInfo(writer *NetLib.RawPacketData, user *RoomUser){

}

func (room *BaseRoom) DisConnectedUser(sessionUniqueID uint64) bool {
	return true
}



func (room *BaseRoom) SecondTimeEvent(){

}

func (room *BaseRoom) BroadCastPacket(){

}


func (room *BaseRoom) DisconnectedUser(sessionUniqueID uint64) int16 {
	return  protocol.ERROR_CODE_NONE
}