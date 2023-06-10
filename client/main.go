package main

import (
	"fmt"
	"github.com/byteYuFan/NAT/instance"
	"github.com/byteYuFan/NAT/network"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

// 按照我们的开发流程，我们需要定义些许常量
// GOOS=linux GOARCH=amd64 go build -o server

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
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		// 只打印帮助信息，不执行命令
		rootCmd.SetArgs(os.Args[1:])
		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		exchange()
		art()
		printRelationInformation()
		// 连接服务器控制接口
		controllerTCPConn, err := network.CreateTCPConn(objectConfig.ControllerAddr)
		if err != nil {
			log.Println("[CreateTCPConn]" + objectConfig.ControllerAddr + err.Error())
			return
		}
		fmt.Println("[Conn Successfully]" + objectConfig.ControllerAddr)
		err = authTheServer(controllerTCPConn)
		if err != nil {
			fmt.Println("[authTheServer]", err)
			return
		}
		nsi := instance.NewSendAndReceiveInstance(controllerTCPConn)
	receiveLoop:
		for {
			msg, err := nsi.ReadHeadDataFromClient()
			if err == io.EOF {
				break
			}
			if opErr, ok := err.(*net.OpError); ok {
				fmt.Println("[err]", err)
				if strings.Contains(opErr.Error(), "An existing connection was forcibly closed by the remote host") {
					// 远程主机关闭连接，退出连接处理循环
					fmt.Println("远程主机关闭连接")
					break
				}
			}
			if err != nil {
				fmt.Println("[err]", err)
				continue
			}
			msg, err = nsi.ReadRealDataFromClient(msg)
			if err != nil {
				fmt.Println("[readReal]", msg)
				continue
			}
			if msg.GetMsgID() == network.AUTH_FAIL {
				fmt.Println("[auth  fail]", "认证失败")
				break
			}
			switch msg.GetMsgID() {
			case network.NEW_CONNECTION:
				processNewConnection(msg.GetMsgData())
			case network.USER_INFORMATION:
				err := processUserInfo(msg.GetMsgData())
				if err != nil {
					fmt.Println("[User Info]", err)
					continue
				}
			case network.KEEP_ALIVE:
				processKeepLive(msg.GetMsgData())
			case network.CONNECTION_IF_FULL:
				processConnIsFull(msg.GetMsgData())
				break receiveLoop
			}

		}
	}
	fmt.Println("[客户端退出，欢迎您的使用。]", "GoodBye, Have a good time!!!")
}

func init() {
	initConfig()
	initCobra()
}
