package main

import (
	"bufio"
	"fmt"
	"github.com/byteYuFan/NAT/network"
	"io"
	"log"
	"net"
)

// 按照我们的开发流程，我们需要定义些许常量

var (
	// 公网ip
	publicAddr = "8.8.8.8"
	// 本地服务的地址
	localServerAddr = "127.0.0.1:8080"
	// 公网服务端的控制接口
	controllerServerAddr = publicAddr + ":8080"
	// 公网隧道地址
	tunnelServerAddr = publicAddr + ":8008"
)

func main() {
	// 与服务器的控制接口建立TCP连接 使用我们工具包的函数
	controllerTCPConn, err := network.CreateTCPConn(controllerServerAddr)
	if err != nil {
		log.Println("[CreateTCPConn]" + controllerServerAddr + err.Error())
		return
	}
	log.Println("[Conn Successfully]" + controllerServerAddr)
	// 新建一个Reader从控制通道中进行连接
	reader := bufio.NewReader(controllerTCPConn)
	// 不断的读取从通道读取信息
	for {
		line, err := reader.ReadString('\n')
		if err != nil || err == io.EOF {
			log.Println("[Controller ReadSting]" + err.Error())
			break
		}
		// 接收到连接的信号
		if line == network.NewConnection+"\n" {
			// 创建连接
			fmt.Println("[创建管道]")
			go connectLocalAndTunnel()
		}

	}
}

// 连接隧道和本地服务器
func connectLocalAndTunnel() {
	local := connLocalServer()
	tunnel := connWebServer()
	network.SwapConnDataEachOther(local, tunnel)
}

// 连接本地web服务器
func connLocalServer() *net.TCPConn {
	local, err := network.CreateTCPConn(localServerAddr)
	if err != nil {
		log.Println("[CreateLocalServerConn]" + err.Error())
		panic(err)
	}
	log.Println("[连接本地服务器成功]")
	return local
}

// 连接web服务器隧道
func connWebServer() *net.TCPConn {
	tunnel, err := network.CreateTCPConn(tunnelServerAddr)
	if err != nil {
		log.Println("[CreateTunnelServerConn]" + err.Error())
		panic(err)
	}
	log.Println("[连接服务器隧道成功]")
	return tunnel
}
