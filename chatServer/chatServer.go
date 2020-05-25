package main

import (


	NetLib "gohipernetFake"
)

type configAppServer struct {
	RoomMaxCount		int32
	RoomStartNum		int32
	RoomMaxUserCount	int32
}

type chatServer struct {
	ServerIndex int
	IP			string
	Port		int

	//PacketChannel	chan protocol
}


func StartServer(netConfig NetLib.NetworkConfig){
	sessionNetworkFunctor := NetLib.SessionNetworkFunctors{}
	//sessionNetworkFunctor.OnReceiveBufferedData



	NetLib.NetLibStartNetwork(&netConfig, sessionNetworkFunctor)
}

func OnReceiveBufferedData(sessionIndex int32, sessionID uint64,data []byte) bool {



	return true;
}


func ProcessPacketData(){

}