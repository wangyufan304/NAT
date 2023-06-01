package main

import (
	"github.com/byteYuFan/NAT/network"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

// 这个版本只实现最基本的NAT穿透，即就是最简单的转发
// 流程大概如下

const (
	// 控制信息地址
	controllerAddr = "0.0.0.0:8009"
	// 隧道地址
	tunnelAddr = "0.0.0.0:8008"
	// 外部访问地址
	visitAddr = "0.0.0.0:8007"
)

var (
	// 内网客户端连接 目前只支持一个客户端进行连接
	clientConn *net.TCPConn
	// 用户连接池
	connectionPool map[string]*UserConnInfo
	// 保护连接池用的锁
	connectionPoolLock sync.Mutex
)

// UserConnInfo 用户连接信息,此处保存的是用户访问公网web所对应的那个接口
type UserConnInfo struct {
	// visit 用户访问web服务的时间
	visit time.Time
	// conn tcp连接句柄
	conn *net.TCPConn
}

func main() {
	// 初始化连接池默认大小为128
	connectionPool = make(map[string]*UserConnInfo, 128)
	go createControllerChannel()
	go acceptUserRequest()
	go acceptClientRequest()
	cleanExpireConnPool()
}

// createControllerChannel 创建一个控制信息的通道，用于传递控制消息
func createControllerChannel() {
	controllerListener, err := network.CreateTCPListener(controllerAddr)
	if err != nil {
		log.Println("[CreateControllerTCPConn]" + controllerAddr + err.Error())
		panic(err)
	}
	log.Println("[公网服务器控制端开始监听]" + controllerAddr)
	for {
		tcpConn, err := controllerListener.AcceptTCP()
		if err != nil {
			log.Println("[ControllerAccept]", err)
			continue
		}
		log.Println("[控制层接收到新的连接]", tcpConn.RemoteAddr())
		// 如果全局变量不为空的话，丢弃该连接
		if clientConn != nil {
			_ = tcpConn.Close()
		} else {
			clientConn = tcpConn
		}

	}
}

// 监听来自用户的请求
func acceptUserRequest() {
	listener, err := network.CreateTCPListener(visitAddr)
	if err != nil {
		log.Println("[CreateVisitListener]" + err.Error())
		panic(err)
	}
	defer listener.Close()
	for {
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println("[VisitListener]", err)
			continue
		}
		addUserConnIntoPool(tcpConn)
		// 向控制通道发送信息
		sendMessageToClientController(network.NewConnection + "\n")
	}
}

// sendMessageToClientController 向客户端发送控制信息
func sendMessageToClientController(message string) {
	if clientConn == nil {
		log.Println("[SendMessage]", "没有连接的客户端")
		return
	}
	_, err := clientConn.Write([]byte(message))
	if err != nil {
		log.Println("[SendMessageWrite]", err)
	}
}

// 接收客户端的请求并建立隧道
func acceptClientRequest() {
	tunnelListener, err := network.CreateTCPListener(tunnelAddr)
	if err != nil {
		log.Println("[CreateTunnelListener]" + tunnelAddr + err.Error())
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
	connectionPoolLock.Lock()
	defer connectionPoolLock.Unlock()

	for key, connMatch := range connectionPool {
		if connMatch.conn != nil {
			go network.SwapConnDataEachOther(connMatch.conn, tunnel)
			delete(connectionPool, key)
			return
		}
	}

	_ = tunnel.Close()
}

// 将用户来的连接放入连接池中
func addUserConnIntoPool(conn *net.TCPConn) {
	connectionPoolLock.Lock()
	defer connectionPoolLock.Unlock()

	now := time.Now()
	connectionPool[strconv.FormatInt(now.UnixNano(), 10)] = &UserConnInfo{now, conn}
}

// cleanExpireConnPool 清理连接池
func cleanExpireConnPool() {
	for {
		connectionPoolLock.Lock()
		for key, connMatch := range connectionPool {
			if time.Now().Sub(connMatch.visit) > time.Second*10 {
				_ = connMatch.conn.Close()
				delete(connectionPool, key)
			}
		}
		connectionPoolLock.Unlock()
		time.Sleep(5 * time.Second)
	}
}
