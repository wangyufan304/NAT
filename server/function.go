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
		// 给客户端发送该消息
		log.Println("[控制层接收到新的连接]", tcpConn.RemoteAddr())
		// 将新的连接推入工作队列中去
		serverInstance.WorkerBuffer <- tcpConn
		fmt.Printf("[%s] %s\n", tcpConn.RemoteAddr().String(), "已推入工作队列中。")
	}
}

// keepAlive 心跳包检测
func keepAlive(conn *net.TCPConn, port int32) {
	for {
		nsi := instance.NewSendAndReceiveInstance(conn)
		_, err := nsi.SendDataToClient(network.KEEP_ALIVE, []byte("ping"))
		if err != nil {
			serverInstance.ProcessWorker.Remove(port)
			return
		}
		time.Sleep(time.Second * 3)
	}
}

// ListenTaskQueue 监听任务队列，获取里面的请求
func ListenTaskQueue() {

	fmt.Println("[ListenTaskQueue]", "监听工作队列传来的消息")
restLabel:
	if !serverInstance.PortIsFull() {
		conn := <-serverInstance.WorkerBuffer
		go acceptUserRequest(conn)
	}
	time.Sleep(time.Millisecond * 10)
	goto restLabel
}

// acceptUserRequest 接收用户的请求
func acceptUserRequest(conn *net.TCPConn) {
	port := serverInstance.GetPort()
	userVisitAddr := "0.0.0.0:" + strconv.Itoa(int(port))
	userVisitListener, err := network.CreateTCPListener(userVisitAddr)
	if err != nil {
		log.Println("[CreateVisitListener]" + err.Error())
		return
	}
	defer userVisitListener.Close()
	workerInstance := NewWorker(userVisitListener, conn, port)
	serverInstance.ProcessWorker.Add(port, workerInstance)
	c := network.NewClientConnInstance(serverInstance.Counter, port)
	ready, _ := c.ToBytes()
	nsi := instance.NewSendAndReceiveInstance(conn)
	go keepAlive(conn, port)
	_, err = nsi.SendDataToClient(network.USER_INFORMATION, ready)
	if err != nil {
		fmt.Println("[Send Client info]", err)
		return
	}
	fmt.Println("[SendClientInfo Successfully]")
	fmt.Println("[addr]", userVisitListener.Addr().String())
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
