package main

import (
	"fmt"
	"github.com/byteYuFan/NAT/instance"
	"github.com/byteYuFan/NAT/network"
	"github.com/byteYuFan/NAT/utils"
	"log"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// 这个版本只实现最基本的NAT穿透，即就是最简单的转发
// 流程大概如下

const (
	// 控制信息地址
	controllerAddr = "0.0.0.0:8080"
	// 隧道地址
	tunnelAddr = "0.0.0.0:8008"
	// 外部访问地址
	visitAddr = "0.0.0.0:8007"
)

var (
	//
	connPool *instance.ConnPool
	// 内网客户端连接 目前只支持一个客户端进行连接
	clientConn *net.TCPConn
	// 用户连接池
	connectionPool map[string]*UserConnInfo
	// 保护连接池用的锁
	connectionPoolLock sync.Mutex
	// 信号
	readyNATConn chan *net.TCPConn
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
	go Accept()
	go acceptClientRequest()
	cleanExpireConnPool()
	select {}
}
func init() {
	connPool = instance.NewConnPool()
	for i := 0; i < int(connPool.MaxTCPConn); i++ {
		connPool.TaskQueue[i] = make(chan *instance.Request, connPool.BufferSize)
	}
	readyNATConn = make(chan *net.TCPConn, 16)
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
		atomic.AddInt32(&connPool.CurrentConnNum, 1)
		fmt.Println("[Receive CONN]", tcpConn.RemoteAddr().String())
		go dealWithControllerInfo(tcpConn)
	}
}
func dealWithControllerInfo(tcpConn *net.TCPConn) {
	uid := connPool.AddTCPConn(tcpConn)
	// 发送信号给readyChannel
	fmt.Println("[向ready推送消息]")
	readyNATConn <- tcpConn
	readChan := make(chan bool)
	go sendControllerDataToClient(uid, tcpConn, readChan)
	go keepAliveDevice(tcpConn, uid, readChan)
}

func sendControllerDataToClient(uid int64, conn *net.TCPConn, signalChan chan bool) {
	dataControllerInfo := &network.ControllerInfo{
		ID:             uint64(connPool.Counter),
		Port:           connPool.Port[uid],
		CurrentConnNum: uint32(connPool.CurrentConnNum),
	}
	byteData, err := utils.ObjectToBufferStream(dataControllerInfo)
	if err != nil {
		fmt.Println("[TO Buffer]", err)
		return
	}
	// 新建一个message消息的实例
	msgReady := instance.NewMsgPackage(network.CONTROLLER_INFO, byteData)
	dp := instance.NewDataPackage()
	// 封装数据
	s, err := dp.Pack(msgReady)
	if err != nil {
		fmt.Println("[Pack]", err)
		return
	}
	written, err := conn.Write(s)
	if err != nil {
		fmt.Println("[Written]", err)
		return
	}
	fmt.Println("[WriteData Successfully]", written)
	signalChan <- true
	// 启动监听客户端进程
	// TODO 启动一个心跳检测装置
}

// keepAliveDevice 心跳包检测机制
func keepAliveDevice(tcpConn *net.TCPConn, uid int64, signalChan chan bool) {
	//
	select {
	case <-signalChan:
		break
	}
	fmt.Println("[KeepAlive Running]..........")
	go func(conn *net.TCPConn) {
		var heartCount uint64
		for {
			// 封装TCP心跳包包的数据
			keepData := network.NewKeepAlive(heartCount, "ping")
			keepDataReady, err := utils.ObjectToBufferStream(keepData)
			heartCount++
			if err != nil {
				fmt.Println("[ObjectToBufferStream]", err)
				return
			}
			readyMsg := instance.NewMsgPackage(network.KEEP_ALIVE, keepDataReady)
			dp := instance.NewDataPackage()
			readyStream, err := dp.Pack(readyMsg)
			if err != nil {
				fmt.Println("[Pack]", err)
				return
			}
			_, err = tcpConn.Write(readyStream)
			if err != nil {
				break
			}
			//fmt.Println("[Write successfully]", written)
			time.Sleep(5 * time.Second)
		}
		fmt.Println("[客户端关闭]", "心跳检测停止")
		fmt.Println("[客户端从工作队列出列]")
		req := connPool.RemoveConn(uid, tcpConn)
		fmt.Println("[ID]", req.ConnParticularInfo.ID, "退出")
	}(tcpConn)
}

// Accept 监听来自用户的请求
func Accept() {
	for {
		select {
		case conn := <-readyNATConn:
			// 获取连接端口号
			fmt.Println("[接受到创建tcp连接信号]")
			userVisitAddr := "0.0.0.0:" + strconv.Itoa(int(connPool.ConnInfo[conn]))
			go acceptUserRequest(conn, userVisitAddr)
		}
	}
}
func acceptUserRequest(natClientConn *net.TCPConn, userVisitAddr string) {
	listener, err := network.CreateTCPListener(userVisitAddr)
	fmt.Println("[CreateConn]", userVisitAddr)
	if err != nil {
		log.Println("[CreateVisitListener]" + err.Error())
		return
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
		sendMessageToClientController(natClientConn, network.NewConnection)
	}
}

// sendMessageToClientController 向客户端发送控制信息
func sendMessageToClientController(natClientConn *net.TCPConn, message string) {
	msg := instance.NewMsgPackage(1, []byte(message))
	dp := instance.NewDataPackage()
	stream, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("Send", err)
		return
	}
	written, err := natClientConn.Write(stream)
	if err != nil {
		fmt.Println("[Written]", err)
		return
	}
	fmt.Println("[Send Successfully.]", written)
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
