package connectedSession

import (
	"chatServer/protocol"
	"sync/atomic"
)

type Session struct {
	Index 					int32
	NetworkUniqueID 		uint64

	UserID 					[protocol.MAX_USER_ID_BYTE_LENGTH]byte
	UserIDLength 			int8

	ConnectTimeSecond		int64
	RoomNumber				int32
	RoomNumberOfEntering	int32
}

func (session *Session) Init(index int32){
	session.Index = index
	session.Clear()
}



func (session *Session) ClearUserID(){
	session.UserIDLength = 0
}

func (session *Session) Clear(){
	session.ClearUserID()
	session.SetRoomNumber(0, -1,0)
	session.SetConnectTimeSecond(0,0)
}



func (session *Session) GetIndex() int32{
	return session.Index
}

func (session *Session) GetNetworkUniqueID() uint64{
	return atomic.LoadUint64(&session.NetworkUniqueID)
}

func (session *Session) GetNetworkInfo() (int32, uint64) {
	index := session.GetIndex()
	uniqueID := atomic.LoadUint64(&session.NetworkUniqueID)

	return index, uniqueID
}



func (session *Session) SetUserID(userID []byte){
	session.UserIDLength = int8(len(userID))
	copy(session.UserID[:], userID)
}

func (session *Session) GetUserID() []byte {
	return session.UserID[0:session.UserIDLength]
}

func (session *Session) GetUserIDLength() int8{
	return session.UserIDLength
}



func (session *Session) SetConnectTimeSecond(timeSec int64, uniqueID uint64){
	atomic.StoreInt64(&session.ConnectTimeSecond, timeSec)
	atomic.StoreUint64(&session.NetworkUniqueID, uniqueID)
}

func (session *Session) GetConnectTimeSecond() int64{
	return atomic.LoadInt64(&session.ConnectTimeSecond)
}

func (session *Session) SetUser(sessionUserID uint64, userID []byte, currentTimeSecond int64){
	session.SetUserID(userID)
	session.SetRoomNumber(sessionUserID, -1,  currentTimeSecond )

}

func (session *Session) IsAuthorized () bool{
	if session.UserIDLength > 0 {
		return true
	}
	return false
}


func (session *Session) ValidNetworkUniqueID(uniqueID uint64) bool {
	return atomic.LoadUint64(&session.NetworkUniqueID) == uniqueID
}



func (session *Session) SetRoomEntering(roomNum int32) bool {
	if atomic.CompareAndSwapInt32(&session.RoomNumberOfEntering, -1, roomNum) == false {
		return false
	}
	return true
}

func (session *Session) SetRoomNumber(sessionUniqueID uint64, roomNum int32, curTimeSec int64) bool {
	if roomNum == -1 {

		atomic.StoreInt32(&session.RoomNumber, roomNum)
		atomic.StoreInt32(&session.RoomNumberOfEntering, roomNum)
	}

	if sessionUniqueID != 0 && session.ValidNetworkUniqueID(sessionUniqueID) == false{
		return false
	}

	if atomic.CompareAndSwapInt32(&session.RoomNumber, -1, roomNum) == false {
		return false
	}

	atomic.StoreInt32(&session.RoomNumberOfEntering, roomNum)

	return true
}

func (session *Session) GetRoomNumber() (int32, int32) {
	roomNum := atomic.LoadInt32(&session.RoomNumber)
	roomNumOfEntering := atomic.LoadInt32(&session.RoomNumberOfEntering)

	return roomNum, roomNumOfEntering
}