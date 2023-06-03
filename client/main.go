package main

import (
	"fmt"
	"github.com/byteYuFan/NAT/instance"
	"github.com/byteYuFan/NAT/network"
	"github.com/byteYuFan/NAT/utils"
	"log"
	"net"
	"time"
)

// 按照我们的开发流程，我们需要定义些许常量

var (
	// 公网ip
	publicAddr = "pogf.com.cn"
	// 本地服务的地址
	localServerAddr = "127.0.0.1:8081"
	// 公网服务端的控制接口
	controllerServerAddr = publicAddr + ":8080"
	// 公网隧道地址
	tunnelServerAddr = publicAddr + ":8008"
)
var clientInfo = new(network.ControllerInfo)

var singleChan = make(chan *net.TCPConn)

func main() {
	// 与服务器的控制接口建立TCP连接 使用我们工具包的函数
	controllerTCPConn, err := network.CreateTCPConn(controllerServerAddr)
	if err != nil {
		log.Println("[CreateTCPConn]" + controllerServerAddr + err.Error())
		return
	}
	log.Println("[Conn Successfully]" + controllerServerAddr)
	go dealWithControllerInformation(controllerTCPConn)
	select {
	case conn := <-singleChan:
		conn.Close()
	}
}

// 处理控制连接信息
func dealWithControllerInformation(controllerTCPConn *net.TCPConn) {
	defer controllerTCPConn.Close()
	for {
		if clientInfo.CurrentConnNum > 1 {
			time.Sleep(time.Second)
			fmt.Println("[连接数量过大，出局]")
			continue
		}
		var keep *network.KeepAlive
		time.Sleep(time.Second * 1)
		HeadData, err := utils.ReadHeadData(controllerTCPConn)
		dp := instance.NewDataPackage()
		msg, err := dp.Unpack(HeadData)
		if err != nil {
			fmt.Println("[UNPACK ]", err)
			return
		}
		var data []byte
		if msg.GetMsgDataLen() > 0 {
			data, err = utils.GetRealData(msg.GetMsgDataLen(), controllerTCPConn)
			if err != nil {
				fmt.Println("[GetData]", err)
				continue
			}
		}
		switch msg.GetMsgID() {
		case network.CONTROLLER_INFO:
			err = dealWithControllerData(clientInfo, data)
			if err != nil {
				fmt.Println("[ControllerData]", err)
				return
			}
			fmt.Println(clientInfo)
			// TODO 收到控制信息后建立一系列连接
			go connectLocalAndTunnel()
		case network.KEEP_ALIVE:
			err = dealWithKeepAliveData(&keep, data)
			if err != nil {
				fmt.Println("[dealWithKeepAliveData]", err)
				continue
			}
			fmt.Println("[ID]", keep.ID, "[MSG]", string(keep.Msg))
		case 1:

			fmt.Println("[New]", string(data))
			go connectLocalAndTunnel()
		}

	}
}

// 连接隧道和本地服务器
func connectLocalAndTunnel() {
	// 修改端口信息

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
