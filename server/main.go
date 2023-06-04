package main

import (
	"fmt"
	"github.com/byteYuFan/NAT/network"
	"github.com/spf13/cobra"
	"log"
	"net"
	"os"
	"time"
)

// 这个版本只实现最基本的NAT穿透，即就是最简单的转发
// 流程大概如下

var rootCmd = &cobra.Command{
	Use:   "Server-NAT [OPTIONS] COMMAND",
	Short: "Intranet penetration",
	Long:  "GO language-based Intranet penetration tool that supports multiple connections",
	Run: func(cmd *cobra.Command, args []string) {
		// 运行命令的处理逻辑
	},
}

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
		go createControllerChannel()
		go ListenTaskQueue()
		go acceptClientRequest()
		go cleanExpireConnPool()
		select {}
	}

}

//func main() {
//	// 初始化连接池默认大小为128
//	//connectionPool = make(map[string]*UserConnInfo, 128)
//	//go createControllerChannel()
//	//go acceptUserRequest()
//	//go acceptClientRequest()
//	//cleanExpireConnPool()
//	fmt.Println("[Object]", objectConfig)
//}

// createControllerChannel 创建一个控制信息的通道，用于传递控制消息
func createControllerChannel() {
	controllerListener, err := network.CreateTCPListener(objectConfig.ControllerAddr)
	if err != nil {
		fmt.Println("[createTCPListener]", err)
		panic(err)
	}
	log.Println("[公网服务器控制端开始监听]" + objectConfig.ControllerAddr)
	for {
		tcpConn, err := controllerListener.AcceptTCP()
		if err != nil {
			fmt.Println("[AcceptTCP]", err)
			continue
		}
		log.Println("[控制层接收到新的连接]", tcpConn.RemoteAddr())
		// 将tcpConn加入连接池中去
		err = serverInstance.AddClientConn(tcpConn)
		if err != nil {
			fmt.Println("[AddClientConn]", err)
			continue
		}
		log.Println("[添加client到连接池中去]")
		// 将tcpConn加入对应的工作队列中去
		err = serverInstance.PushConnToTaskQueue(tcpConn)
		if err != nil {
			fmt.Println("[PushConnToTaskQueue]", err)
			continue
		}
		log.Println("[添加client到工作队列去]")
	}
}

// createControllerChannel 创建一个控制信息的通道，用于传递控制消息
//func createControllerChannel() {
//	controllerListener, err := network.CreateTCPListener(controllerAddr)
//	if err != nil {
//		log.Println("[CreateControllerTCPConn]" + controllerAddr + err.Error())
//		panic(err)
//	}
//	log.Println("[公网服务器控制端开始监听]" + controllerAddr)
//	for {
//		tcpConn, err := controllerListener.AcceptTCP()
//		if err != nil {
//			log.Println("[ControllerAccept]", err)
//			continue
//		}
//		log.Println("[控制层接收到新的连接]", tcpConn.RemoteAddr())
//		// 如果全局变量不为空的话，丢弃该连接
//		if clientConn != nil {
//			_ = tcpConn.Close()
//		} else {
//			clientConn = tcpConn
//		}
//
//	}
//}

// ListenTaskQueue 监听任务队列，获取里面的请求
func ListenTaskQueue() {
	fmt.Println("[ListenTaskQueue]", "监听工作队列传来的消息")
	for i := 0; i < int(serverInstance.MaxTCPConnSize); i++ {
		go func(num int) {
			for {
				conn := <-serverInstance.TaskQueueSlice[num].Worker
				go acceptUserRequest(conn)
			}
		}(i)
	}
}

// acceptUserRequest 接收用户的请求
func acceptUserRequest(conn *net.TCPConn) {
	// 根据用户的conn先从全局map中获取到它的uid
	uid, err := serverInstance.GetClientUid(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 根据uid获取到对应的ip地址
	getPort := serverInstance.TaskQueueSlice[uid%int64(serverInstance.TaskQueueSize)].GetPort()
	userVisitAddr := fmt.Sprintf("%s:%d", "0.0.0.0", getPort)
	userVisitListener, err := network.CreateTCPListener(userVisitAddr)
	if err != nil {
		log.Println("[CreateVisitListener]" + err.Error())
		panic(err)
	}
	fmt.Println("[CreateTCPListener successfully]", userVisitAddr)
	defer fmt.Println("[关闭 successfully]", userVisitAddr)
	defer userVisitListener.Close()
	for {
		tcpConn, err := userVisitListener.AcceptTCP()
		if err != nil {
			log.Println("[userVisitListener.AcceptTCP]", err)
			continue
		}
		userConnPoolInstacne.AddConnInfo(tcpConn)
		nsi := NewControllerServiceInstance(conn)
		count, err := nsi.SendDataToClient(1, []byte(network.NewConnection))
		if err != nil {
			fmt.Println("[SendData fail]", err)
			continue
		}
		fmt.Println("[SendData successfully]", count, " bytes")
	}
}

//// 监听来自用户的请求
//func acceptUserRequest() {
//	listener, err := network.CreateTCPListener(visitAddr)
//	if err != nil {
//		log.Println("[CreateVisitListener]" + err.Error())
//		panic(err)
//	}
//	defer listener.Close()
//	for {
//		tcpConn, err := listener.AcceptTCP()
//		if err != nil {
//			log.Println("[VisitListener]", err)
//			continue
//		}
//		addUserConnIntoPool(tcpConn)
//		// 向控制通道发送信息
//		sendMessageToClientController(network.NewConnection + "\n")
//	}
//}

//// sendMessageToClientController 向客户端发送控制信息
//func sendMessageToClientController(message string) {
//	if clientConn == nil {
//		log.Println("[SendMessage]", "没有连接的客户端")
//		return
//	}
//	_, err := clientConn.Write([]byte(message))
//	if err != nil {
//		log.Println("[SendMessageWrite]", err)
//	}
//}

// 接收客户端的请求并建立隧道
func acceptClientRequest() {
	tunnelListener, err := network.CreateTCPListener(objectConfig.TunnelAddr)
	if err != nil {
		log.Println("[CreateTunnelListener]" + objectConfig.TunnelAddr + err.Error())
		return
	}
	defer tunnelListener.Close()
	for {
		tcpConn, err := tunnelListener.AcceptTCP()
		if err != nil {
			log.Println("[TunnelAccept]", err)
			continue
		}
		// 创建隧道
		go createTunnel(tcpConn)
	}
}

// createTunnel 核心，将用户的请求数据转发给tunnel，然后内网客户端在转发到内网服务器
func createTunnel(tunnel *net.TCPConn) {
	userConnPoolInstacne.Mutex.RLock()
	defer userConnPoolInstacne.Mutex.RUnlock()

	for key, connMatch := range userConnPoolInstacne.UserConnectionMap {
		if connMatch.conn != nil {
			go network.SwapConnDataEachOther(connMatch.conn, tunnel)
			delete(userConnPoolInstacne.UserConnectionMap, key)
			return
		}
	}

	_ = tunnel.Close()
}

//// 将用户来的连接放入连接池中
//func addUserConnIntoPool(conn *net.TCPConn) {
//	connectionPoolLock.Lock()
//	defer connectionPoolLock.Unlock()
//
//	now := time.Now()
//	connectionPool[strconv.FormatInt(now.UnixNano(), 10)] = &UserConnInfo{now, conn}
//}

// cleanExpireConnPool 清理连接池
func cleanExpireConnPool() {
	for {
		userConnPoolInstacne.Mutex.Lock()
		for key, connMatch := range userConnPoolInstacne.UserConnectionMap {
			if time.Now().Sub(connMatch.visit) > time.Second*10 {
				_ = connMatch.conn.Close()
				delete(userConnPoolInstacne.UserConnectionMap, key)
			}
		}
		userConnPoolInstacne.Mutex.Unlock()
		time.Sleep(5 * time.Second)
	}
}

func init() {
	initConfig()
	initCobra()
	initServer()
	initUserConnPool()
}
