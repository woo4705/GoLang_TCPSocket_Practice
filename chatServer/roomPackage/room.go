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
	FuncList				[]func(*RoomUser, protocol.Packet) int16

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



func (room *BaseRoom) ISUserCanEnter() bool {
	if room.GetCurrentUserCount() >= room.Config.MaxRoomCount {
		return false
	}
	return true
}


func (room *BaseRoom) SettingPacketFunction() {
	maxFunctionListCount := 16
	room.FuncList = make([]func(*RoomUser, protocol.Packet)int16, 0, maxFunctionListCount)
	room.FuncPacketIDList = make([]int16, maxFunctionListCount)

	//TODO: AddFunction 작업하기. 함수호출 구현법 더 살펴보고 이해해보기

}

func (room *BaseRoom) AddPacketFunction(packetID int16, packetFunc func(*RoomUser, protocol.Packet) int16){
	room.FuncList = append(room.FuncList, packetFunc)
	room.FuncPacketIDList = append(room.FuncPacketIDList, packetID)
}



func (room *BaseRoom) InitUserPool() {
	room.UserPool = &sync.Pool{
		New: func() interface{} {
			user := new(RoomUser)
			return user
		},
	}
}


func (room *BaseRoom) GetUserObject() *RoomUser {
	userObject := room.UserPool.Get().(*RoomUser)
	return userObject
}

func (room *BaseRoom) PutUserObject(user *RoomUser) {
	room.UserPool.Put(user)
}


func (room *BaseRoom) AddUser(userInfo AddRoomUserInfo) (*RoomUser, int16){
	if room.IsRoomFull(){
		return nil, protocol.ERROR_CODE_ENTER_ROOM_USER_FULL
	}

	if room.GetUser(userInfo.netSessionUniqueID) != nil {
		return nil, protocol.ERROR_CODE_ENTER_ROOM_DUPLICATION_USER
	}

	atomic.AddInt32(&room.CurUserCount, 1)

	user := room.GetUserObject()
	user.Init(userInfo.userID, room.GenerateUserUniqueID())
	user.SetNetworkInfo(userInfo.netSessionIndex, userInfo.netSessionUniqueID)
	user.packetDataSize = user.PacketDataSize()

	room.UserSessionUniqueIDMap[user.netSessionUniqueID] = user
	return user, protocol.ERROR_CODE_NONE
}



func (room *BaseRoom) IsRoomFull() bool {
	if room.GetCurrentUserCount() == room.Config.MaxRoomCount{
		return true
	}
	return false
}



func (room *BaseRoom) RemoveUser(user *RoomUser) {
	delete(room.UserSessionUniqueIDMap, user.NetSessionUniqueID)
	room.RemoveUserObject(user)
}

func (room *BaseRoom) RemoveUserObject(user *RoomUser) {
	atomic.AddInt32(&room.CurUserCount, -1)
	room.PutUserObject(user)
}

func (room *BaseRoom) GetUser(sessionUniqueID uint64) *RoomUser {
	if user, ok := room.UserSessionUniqueIDMap[sessionUniqueID]; ok {
		return user
	}
	return nil
}



func (room *BaseRoom) AllocAllUserInfo(exceptSessionUniqueID uint64) (userCount int8, dataSize int16, dataBuffer []byte) {
	for _, user := range room.UserSessionUniqueIDMap {
		if user.NetSessionUniqueID == exceptSessionUniqueID {
			continue
		}

		userCount++
		dataSize += user.packetDataSize
	}

	dataBuffer = make([]byte, dataSize)
	writer := NetLib.MakeWriter(dataBuffer, true)


	for _, user := range  room.UserSessionUniqueIDMap {
		if user.NetSessionUniqueID == exceptSessionUniqueID {
			continue
		}
		WriteUserInfo(&writer, user)

	}

	return userCount, dataSize, dataBuffer
}


func  (room *BaseRoom)AllocUserInfo(user *RoomUser) (dataSize int16, dataBuffer []byte){
	dataSize = user.packetDataSize
	dataBuffer = make([]byte, dataSize)

	writer := NetLib.MakeWriter(dataBuffer, true)
	WriteUserInfo(&writer, user)

	return dataSize, dataBuffer
}



func WriteUserInfo(writer *NetLib.RawPacketData, user *RoomUser){
	writer.WriteU64(user.RoomUniqueID)
	writer.WriteS8(user.IDLen)
	writer.WriteBytes(user.ID[0:user.IDLen])
}





func (room *BaseRoom) BroadCastPacket(packetSize int16, sendPacket []byte, exceptSessionUniqueID uint64) {
	for _,user := range room.UserSessionUniqueIDMap {
		if user.NetSessionUniqueID == exceptSessionUniqueID {
			continue
		}

		NetLib.NetLibIPostSendToClient(user.NetSessionIndex, user.NetSessionUniqueID, sendPacket)
	}

}


func (room *BaseRoom) IsDisconnectedUser(sessionUniqueID uint64) int16 {
	user := room.GetUser(sessionUniqueID)
	var isUserExist bool

	if user == nil {
		isUserExist = false
	}else {
		isUserExist = true
	}

	if isUserExist == false {
		return protocol.ERROR_CODE_LEAVE_ROOM_INTERNAL_INVALID_USER
	}

	return  protocol.ERROR_CODE_NONE
}