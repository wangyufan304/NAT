package main

import (
	"fmt"
	"github.com/byteYuFan/NAT/network"
)

func processNewConnection(data []byte) {
	// TODO 目前从服务端发送来的信息没有进行处理，后续考虑进行处理
	// 发送自身信息
	go connectLocalAndTunnel()
}

func processUserInfo(data []byte) error {
	ci := network.NewClientConnInstance(0, 0)
	if err := ci.FromBytes(data); err != nil {
		return err
	}
	clientInfo = ci
	return nil
}

func processKeepLive(data []byte) {
	// TODO 目前只简简单单接收服务端发来的请求，简单的打印一下
	fmt.Println("[receive KeepLive package]", string(data))
}

func processConnIsFull(data []byte) {
	fmt.Println("[Server Error]", string(data))
}
