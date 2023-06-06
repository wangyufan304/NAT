package main

import (
	"fmt"
	"github.com/byteYuFan/NAT/instance"
	"github.com/byteYuFan/NAT/network"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

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
		if serverInstance.GetNumOfCurrentConn() >= int(serverInstance.MaxTCPConnSize) {
			nsi := instance.NewSendAndReceiveInstance(tcpConn)
			_, _ = nsi.SendDataToClient(network.CONNECTION_IF_FULL, []byte(network.ProtocolMap[network.CONNECTION_IF_FULL].(string)))
			fmt.Println("[SendFull]")
			continue
		}
		// 给客户端发送该消息
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
		go keepAlive(tcpConn)
	}
}

// keepAlive 心跳包检测
func keepAlive(conn *net.TCPConn) {
	for {
		nsi := instance.NewSendAndReceiveInstance(conn)
		_, err := nsi.SendDataToClient(network.KEEP_ALIVE, []byte("ping"))
		if err != nil {
			// 关闭对应的服务端连接
			serverInstance.Mutex.RLock()
			fmt.Println("[关闭用户访问端口]", serverInstance.ListenerAndClientConn[conn].Addr().String())
			serverInstance.Mutex.RUnlock()
			// 关闭服务端的端口
			serverInstance.ListenerAndClientConn[conn].Close()
			serverInstance.RemoveClientConn(conn)
			return
		}
		time.Sleep(time.Second * 3)
	}
}

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
	cci := network.NewClientConnInstance(serverInstance.Counter, int32(serverInstance.GetConnPortByUID(uid)))
	cciStream, err := cci.ToBytes()
	if err != nil {
		fmt.Println("[ToBytes]", err)
		return
	}
	nsi := instance.NewSendAndReceiveInstance(conn)
	_, err = nsi.SendDataToClient(network.USER_INFORMATION, cciStream)
	if err != nil {
		fmt.Println("[Send UserInfo]", err)
		return
	}
	fmt.Println("[SendClientInfo Successfully]")
	// 根据uid获取到对应的ip地址
	getPort := serverInstance.TaskQueueSlice[uid%int64(serverInstance.TaskQueueSize)].GetPort()
	fmt.Println("[Port]", uid, getPort)
	userVisitAddr := "0.0.0.0:" + strconv.Itoa(getPort)
	userVisitListener, err := network.CreateTCPListener(userVisitAddr)
	if err != nil {
		log.Println("[CreateVisitListener]" + err.Error())

		return
	}
	fmt.Println("[addr]", userVisitListener.Addr().String())
	serverInstance.AddListenerAndClient(userVisitListener, conn)
	fmt.Println("[CreateTCPListener successfully]", userVisitAddr)
	defer fmt.Println("[关闭 successfully]", userVisitAddr)
	defer userVisitListener.Close()
	for {
		tcpConn, err := userVisitListener.AcceptTCP()
		if opErr, ok := err.(*net.OpError); ok {
			if strings.Contains(opErr.Error(), "use of closed network connection") {
				// 远程主机关闭连接，退出连接处理循环
				fmt.Println("远程客户端连接关闭")
				return
			}
		}
		if err != nil {
			log.Println("[userVisitListener.AcceptTCP]", err)
			continue
		}
		userConnPoolInstacne.AddConnInfo(tcpConn)
		nsi := instance.NewSendAndReceiveInstance(conn)
		count, err := nsi.SendDataToClient(1, []byte(network.NewConnection))
		if err != nil {
			fmt.Println("[SendData fail]", err)
			continue
		}
		fmt.Println("[SendData successfully]", count, " bytes")
	}
}

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
