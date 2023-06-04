package main

import (
	"fmt"
	"github.com/byteYuFan/NAT/network"
	"github.com/spf13/cobra"
	"log"
	"net"
	"strings"
)

// 按照我们的开发流程，我们需要定义些许常量
// GOOS=linux GOARCH=amd64 go build -o output_filename

var rootCmd = &cobra.Command{
	Use:   "Client [OPTIONS] COMMAND",
	Short: "Intranet penetration client program.",
	Long:  "If the intranet is written in the go language, you need to start the intranet client before you can connect",
	Run: func(cmd *cobra.Command, args []string) {
		// 运行命令的处理逻辑
	},
}
var clientInfo *network.ClientConnInfo

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("[Cobra]", err)
		panic(err)
	}
	// 与服务器的控制接口建立TCP连接 使用我们工具包的函数
	fmt.Println(objectConfig)
	controllerTCPConn, err := network.CreateTCPConn(objectConfig.ControllerAddr)
	if err != nil {
		log.Println("[CreateTCPConn]" + objectConfig.ControllerAddr + err.Error())
		return
	}
	log.Println("[Conn Successfully]" + objectConfig.ControllerAddr)
	nsi := NewControllerServiceInstance(controllerTCPConn)
	for {
		msg, err := nsi.ReadHeadDataFromClient()
		if opErr, ok := err.(*net.OpError); ok {
			if strings.Contains(opErr.Error(), "An existing connection was forcibly closed by the remote host") {
				// 远程主机关闭连接，退出连接处理循环
				fmt.Println("远程主机关闭连接")
				break
			}
		}
		if err != nil {
			continue
		}
		msg, err = nsi.ReadRealDataFromClient(msg)
		if err != nil {
			fmt.Println("[readReal]", msg)
			continue
		}
		if string(msg.GetMsgData()) == network.NewConnection {
			//创建连接
			fmt.Println("[创建管道]")
			go connectLocalAndTunnel()
		}
		if msg.GetMsgID() == network.USER_INFORMATION {
			ci := network.NewClientConnInstance(0, 0)
			fmt.Println("[Byte]", msg.GetMsgData())
			if err := ci.FromBytes(msg.GetMsgData()); err != nil {
				fmt.Println("[GetMsgData]", err)
				continue
			}
			clientInfo = ci
			fmt.Println("[ClientInfo]", ci)
			fmt.Println("[ClientInfo]", clientInfo)
		}
		if msg.GetMsgID() == network.KEEP_ALIVE {
			fmt.Println("[Heart]", string(msg.GetMsgData()))
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

func init() {
	initConfig()
	initCobra()
}
