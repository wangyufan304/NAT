package main

import (
	"fmt"
	"github.com/byteYuFan/NAT/instance"
	"github.com/byteYuFan/NAT/network"
	"net"
	"strconv"
	"strings"
	"time"
)

// createControllerChannel 创建一个控制信息的通道，用于接收内网客户端的连接请求
// 当内网客户端向服务端的控制接口发送请求建立连接时，控制端会直接向全局的工作队列中添加这个连接信息
// 可以在此进行用户权限的界别与控制
// Create a control information channel to receive connection requests from intranet clients;
// When an intranet client sends a connection request to the control interface of the server,
// the control side will directly add this connection information to the global work queue.
// You can implement user-level permissions and control at this point.
func createControllerChannel() {
	controllerListener, err := network.CreateTCPListener(objectConfig.ControllerAddr)
	if err != nil {
		fmt.Println("[createTCPListener]", err)
		panic(err)
	}
	fmt.Println("[服务器控制端开始监听]" + objectConfig.ControllerAddr)
	if objectConfig.StartAuth == "true" {
		// 获取用户发送来的数据
		fmt.Println("[Start Auth Successfully!]", "服务器开启认证请求")
	}
	for {
		tcpConn, err := controllerListener.AcceptTCP()
		if err != nil {
			fmt.Println("[AcceptTCP]", err)
			continue
		}
		if objectConfig.StartAuth == "true" {
			err := authUser(tcpConn)
			if err != nil {
				fmt.Println(err)
				usi := instance.NewSendAndReceiveInstance(tcpConn)
				_, err = usi.SendDataToClient(network.AUTH_FAIL, []byte{})
				fmt.Println("[AUTH_FAIL]", "发送认证失败消息")
				_ = tcpConn.Close()
				continue
			}
		}
		// 给客户端发送该消息
		fmt.Println("[控制层接收到新的连接]", tcpConn.RemoteAddr())
		// 将新地连接推入工作队列中去
		serverInstance.WorkerBuffer <- tcpConn
		log.Infoln("[%s] %s\n", tcpConn.RemoteAddr().String(), "已推入工作队列中。")
	}
}

// keepAlive 心跳包检测,函数负责向客户端发送保活消息以确保连接处于活动状态(每三秒发送一次)。如果在此过程中发生错误，它会检查错误是否表示客户端已关闭连接。
// 如果是，则会记录相应的日志，并从工作队列中移除相应的端口。然后函数返回。
// The keepAlive function is responsible for sending a keep-alive message to the client to ensure the connection is active.
// If an error occurs during the process, it checks if the error indicates that the client has closed the connection.
// If so, it logs the appropriate message and removes the corresponding port from the worker queue.
// The function then returns.
func keepAlive(conn *net.TCPConn, port int32) {
	for {
		nsi := instance.NewSendAndReceiveInstance(conn)
		_, err := nsi.SendDataToClient(network.KEEP_ALIVE, []byte("ping"))
		log.Infoln("SendData [ping] Successfully.")
		if err != nil {
			log.Errorln("[检测到客户端关闭]", err)
			serverInstance.ProcessWorker.Remove(port)
		}
		time.Sleep(time.Minute)
	}
}

// ListenTaskQueue 该函数的作用是监听工作队列传来的消息。
// 它通过不断检查工作队列是否有可用的连接，并将连接分配给处理函数 acceptUserRequest。
// 当工作队列未满时，会从工作队列中取出一个连接，并启动一个协程来处理该连接的用户请求。
// 函数会以很小的时间间隔进行轮询，并持续监听工作队列的新消息。
// The function listens for messages from the work queue.
// It does this by constantly checking the work queue for an available connection and assigning the connection to the handler function acceptUserRequest.
// When the work queue is not full, a connection is taken from the work queue and a coroutine is started to process user requests for that connection.
// The function polls at small intervals and continuously listens for new messages from the work queue.
func ListenTaskQueue() {
	log.Infoln("[ListenTaskQueue]", "监听工作队列传来的消息")
restLabel:
	if !serverInstance.PortIsFull() {
		conn := <-serverInstance.WorkerBuffer
		go acceptUserRequest(conn)
	}
	time.Sleep(time.Millisecond * 100)
	goto restLabel
}

// acceptUserRequest 接收请用户的求,该函数会首先从全局工作池中获取一个空闲的端口，然后在这个端口上监听用户的请求
// 并向客户端发送对应的信息，然后在这个端口监听用户的请求，每监听到一个请求，就向内网客户端发送一个建立通道的信号
// The function acceptUserRequest is responsible for accepting user requests.
// It first retrieves an available port from the global worker pool.
// It then listens for incoming requests on this port and sends corresponding information to the client.
// Each time a request is received, it sends a signal to establish a channel with the internal network client.
func acceptUserRequest(conn *net.TCPConn) {
	port := serverInstance.GetPort()
	userVisitAddr := "0.0.0.0:" + strconv.Itoa(int(port))
	userVisitListener, err := network.CreateTCPListener(userVisitAddr)
	if err != nil {
		log.Errorln("[CreateTCPListener]", err)
		return
	}
	defer userVisitListener.Close()
	workerInstance := NewWorker(userVisitListener, conn, port)
	serverInstance.ProcessWorker.Add(port, workerInstance)
	c := network.NewClientConnInstance(serverInstance.Counter, port)
	ready, _ := c.ToBytes()
	nsi := instance.NewSendAndReceiveInstance(conn)
	go keepAlive(conn, port)
	_, err = nsi.SendDataToClient(network.USER_AUTHENTICATION_SUCCESSFULLY, []byte{})
	_, err = nsi.SendDataToClient(network.USER_INFORMATION, ready)
	if err != nil {
		log.Infoln("[Send Client info]", err)
		return
	}
	log.Infoln("[addr]", userVisitListener.Addr().String())
	for {
		tcpConn, err := userVisitListener.AcceptTCP()
		if opErr, ok := err.(*net.OpError); ok {
			if strings.Contains(opErr.Error(), "use of closed network connection") {
				// 远程主机关闭连接，退出连接处理循环
				log.Infoln("远程客户端连接关闭")
				return
			}
		}
		if err != nil {
			log.Errorln("[userVisitListener.AcceptTCP]", err)
			continue
		}
		userConnPoolInstance.AddConnInfo(tcpConn)
		nsi := instance.NewSendAndReceiveInstance(conn)
		count, err := nsi.SendDataToClient(network.NEW_CONNECTION, []byte(network.NewConnection))
		if err != nil {
			log.Errorln("[SendData fail]", err)
			continue
		}
		log.Infoln("[SendData successfully]", count, " bytes")
	}
}

// acceptClientRequest 该函数用于接收客户端的请求连接。它首先创建一个监听指定隧道地址的TCP监听器。如果创建监听器时发生错误，则记录错误并返回。函数执行完毕后会关闭监听器。
// 然后，函数进入一个无限循环，接受来自客户端的TCP连接请求。每当接收到一个连接请求时，会创建一个新的协程来处理该连接，即创建隧道。
// This function is responsible for accepting client connection requests.
// It first creates a TCP listener for the specified tunnel address.
// If an error occurs during the creation of the listener, the error is logged and the function returns.
// The listener is closed when the function completes.
// Then, the function enters an infinite loop to accept TCP connection requests from clients. For each incoming connection request, a new goroutine is spawned to handle the connection by creating a tunnel.
func acceptClientRequest() {
	tunnelListener, err := network.CreateTCPListener(objectConfig.TunnelAddr)
	if err != nil {
		log.Errorln("[CreateTunnelListener]" + objectConfig.TunnelAddr + err.Error())
		return
	}
	defer tunnelListener.Close()
	for {
		tcpConn, err := tunnelListener.AcceptTCP()
		if err != nil {
			log.Errorln("[TunnelAccept]", err)
			continue
		}
		// 创建隧道
		go createTunnel(tcpConn)
	}
}

// createTunnel 该函数用于创建一个隧道。
// 函数首先获取用户连接池的读锁，以保证在创建隧道期间不会有其他线程修改连接池。
// 然后它遍历用户连接池中的每个连接，找到一个可用的连接，将该连接与传入的隧道进行数据交换，然后从连接池中删除该连接。
// 如果没有找到可用连接，函数会关闭传入的隧道。最后，释放用户连接池的读锁。
// This function is used to create a tunnel.
// It first acquires the read lock of the user connection pool to ensure that no other threads modify the connection pool during the creation of the tunnel.
// Then it iterates through each connection in the user connection pool to find an available connection.
// It swaps data between the found connection and the provided tunnel, and then removes the connection from the connection pool.
// If no available connection is found, the function closes the provided tunnel. Finally, it releases the read lock of the user connection pool.
func createTunnel(tunnel *net.TCPConn) {
	userConnPoolInstance.Mutex.RLock()
	defer userConnPoolInstance.Mutex.RUnlock()

	for key, connMatch := range userConnPoolInstance.UserConnectionMap {
		if connMatch.conn != nil {
			go network.SwapConnDataEachOther(connMatch.conn, tunnel)
			delete(userConnPoolInstance.UserConnectionMap, key)
			return
		}
	}

	_ = tunnel.Close()
}

// cleanExpireConnPool 该函数用于清理连接池中的过期连接。
// 函数会进入一个无限循环，在每次循环中，它获取连接池的互斥锁，遍历连接池中的每个连接。
// 如果某个连接的访问时间距离当前时间已经超过10秒，那么该连接会被关闭，并从连接池中删除。
// 完成遍历后，释放连接池的互斥锁。函数会每隔5秒执行一次清理操作。
// This function is responsible for cleaning up expired connections in the connection pool.
// It enters an infinite loop, and in each iteration, it acquires the mutex lock of the connection pool.
// It iterates through each connection in the pool.
// If the time elapsed since the last visit of a connection exceeds 10 seconds, the connection is closed and removed from the connection pool.
// After the iteration is complete, the mutex lock of the connection pool is released.
// The function performs the cleanup operation every 5 seconds.
func cleanExpireConnPool() {
	for {
		userConnPoolInstance.Mutex.Lock()
		for key, connMatch := range userConnPoolInstance.UserConnectionMap {
			if time.Now().Sub(connMatch.visit) > time.Second*8 {
				_ = connMatch.conn.Close()
				delete(userConnPoolInstance.UserConnectionMap, key)
			}
		}
		log.Infoln("[cleanExpireConnPool successfully]")
		userConnPoolInstance.Mutex.Unlock()
		time.Sleep(5 * time.Second)
	}
}
