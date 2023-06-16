package main

import (
	"fmt"
	"github.com/byteYuFan/NAT/instance"
	"github.com/byteYuFan/NAT/network"
	"log"
	"net"
)

// 连接隧道和本地服务器
func connectLocalAndTunnel() {
	local := connLocalServer()
	tunnel := connWebServer()
	nsi := instance.NewSendAndReceiveInstance(tunnel)
	stream, err := clientInfo.ToBytes()
	if err != nil {
		fmt.Println(err)
		return
	}
	nsi.SendDataToClient(network.USER_INFORMATION, stream)
	fmt.Println("发送信息successfully")
	network.SwapConnDataEachOther(local, tunnel)
}

// 连接本地web服务器
func connLocalServer() *net.TCPConn {
	local, err := network.CreateTCPConn(objectConfig.LocalServerAddr)
	if err != nil {
		log.Println("[CreateLocalServerConn]" + err.Error())
		panic(err)
	}
	log.Println("[连接本地服务器成功]")
	return local
}

// 连接web服务器隧道
func connWebServer() *net.TCPConn {
	tunnel, err := network.CreateTCPConn(objectConfig.TunnelServerAddr)
	if err != nil {
		log.Println("[CreateTunnelServerConn]" + err.Error())
		panic(err)
	}
	log.Println("[连接服务器隧道成功]")
	return tunnel
}

// authTheServer 向服务器发送认证消息
func authTheServer(conn *net.TCPConn) error {
	// 新建一个数据结构体
	ui := network.NewUserInfoInstance(0, objectConfig.UserName, objectConfig.Password)
	byteStream, err := ui.ToBytes()
	if err != nil {
		return err
	}
	nsi := instance.NewSendAndReceiveInstance(conn)
	_, err = nsi.SendDataToClient(network.USER_REQUEST_AUTH, byteStream)
	if err != nil {
		return err
	}
	return nil
}
