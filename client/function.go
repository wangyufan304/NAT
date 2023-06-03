package main

import (
	"github.com/byteYuFan/NAT/network"
)

func dealWithControllerData(controller **network.ControllerInfo, data []byte) error {
	*controller = new(network.ControllerInfo)
	return (*controller).BufferToObject(data)
}

func dealWithKeepAliveData(keep **network.KeepAlive, data []byte) error {
	*keep = new(network.KeepAlive)
	(*keep).Msg = make([]byte, 4)
	return (*keep).BufferToObject(data)
}
