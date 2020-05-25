package main

import (
	NetLib "gohipernetFake"
)

func StartServer(netConfig NetLib.NetworkConfig){
	sessionNetworkFunctor := NetLib.SessionNetworkFunctors{}
	//sessionNetworkFunctor.OnReceiveBufferedData



	NetLib.NetLibStartNetwork(&netConfig, sessionNetworkFunctor)
}

func OnReceiveBufferedData(sessionIndex int32, sessionID uint64,data []byte) bool {


	return true;
}