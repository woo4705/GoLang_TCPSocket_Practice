package protocol

import (
	"encoding/binary"
	"reflect"
	NetLib "gohipernetFake"
)


const (
	PACKET_TYPE_NORMAL	 = 0
	PACKET_TYPE_COMPRESS = 1
	PACKET_TYPE_SECURE	 = 2
)

const (
	MAX_USER_ID_BYTE_LENGTH	= 16
	MAX_USER_PW_BYTE_LENGTH	= 16
	MAX_CHAT_MESSAGE_BYTE_LENGTH = 126
)

type Header struct {
	TotalSize	int16
	ID			int16
	PacketType	int8
}

type Packet struct {
	UserSessionIndex	int32
	UserSessionUniqueID	uint64
	ID					int16
	DataSize			int16
	Data				[]byte
}


func (packet Packet) GetSessionInfo() (int32,uint64){
	return packet.UserSessionIndex, packet.UserSessionUniqueID
}

var public_clientSessionHeaderSize int16
var public_serverSessionHeaderSize int16

func Init_packetSize() {
	public_clientSessionHeaderSize = ProtocolInitHeaderSize()
	public_serverSessionHeaderSize = ProtocolInitHeaderSize()
}

func ClientHeaderSize() int16{
	return public_clientSessionHeaderSize
}

func ServerHeaderSize() int16{
	return public_serverSessionHeaderSize
}


func ProtocolInitHeaderSize() int16 {
	header := Header{}
	return int16( NetLib.Sizeof(reflect.TypeOf(header)) )
}


func PeekPacketID(rawData []byte) int16{
	packetID := binary.LittleEndian.Uint16(rawData[2:])
	return int16(packetID)
}

func PeekPacketBody(rawData []byte) (bodySize int16, refBody []byte){

	headerSize := ClientHeaderSize()
	bodyData := rawData[:headerSize]
	bodySize = int16( binary.LittleEndian.Uint16(rawData)) - headerSize

	return bodySize, bodyData
}

func DecodingPacketHeader(header *Header, data []byte)  {
	reader := NetLib.MakeReader(data,true)
	header.TotalSize, _ = reader.ReadS16()
	header.ID, _ = reader.ReadS16()
	header.PacketType,_ = reader.ReadS8()
}

func EncodingPacketHeader(writer *NetLib.RawPacketData, totalSize int16, packetID int16, packetType int8)  {
	writer.WriteS16(totalSize)
	writer.WriteS16(packetID)
	writer.WriteS8(packetType)
}
