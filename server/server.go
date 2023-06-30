package main

import (
	"net"
	"sync"
)

var serverInstance *Server

// Server 服务端程序的实例
type Server struct {
	// Mutex 保证并发安全的锁
	Mutex sync.RWMutex
	// Counter 目前服务器累计接收到了多少次连接
	Counter int64
	// 最大连接数量
	MaxTCPConnSize int32
	// 最大连接数量
	MaxConnSize int32
	// ExposePort 服务端暴露端口
	ExposePort []int
	// ProcessingMap
	ProcessingMap map[string]*net.TCPConn
	// WorkerBuffer 整体工作队列的大小
	WorkerBuffer chan *Request
	// 实际处理工作的数据结构
	ProcessWorker *Workers
	// 端口使用情况
	PortStatus      map[int32]bool
	PortStatusMutex sync.RWMutex
	// ConnPort
	ConnPortMap map[int64]int32
	// 新增客户端信息模块
	ClientInfoMap   map[int64]*ClientInfo
	ClientInfoMutex sync.RWMutex
}

func initServer() {
	serverInstance = &Server{
		Mutex:           sync.RWMutex{},
		Counter:         0,
		MaxTCPConnSize:  objectConfig.MaxTCPConnNum,
		MaxConnSize:     objectConfig.MaxConnNum,
		ExposePort:      objectConfig.ExposePort,
		ProcessingMap:   make(map[string]*net.TCPConn),
		WorkerBuffer:    make(chan *Request, objectConfig.MaxConnNum),
		ProcessWorker:   NewWorkers(),
		PortStatus:      make(map[int32]bool),
		PortStatusMutex: sync.RWMutex{},
		ConnPortMap:     make(map[int64]int32),
		ClientInfoMap:   make(map[int64]*ClientInfo),
		ClientInfoMutex: sync.RWMutex{},
	}

	// 初始化端口状态
	for i := 0; i < int(serverInstance.MaxTCPConnSize); i++ {
		serverInstance.PortStatus[int32(serverInstance.ExposePort[i])] = false
	}
}

func (s *Server) PortIsFull() bool {
	s.PortStatusMutex.RLock()
	defer s.PortStatusMutex.RUnlock()
	for _, v := range s.PortStatus {
		if v == false {
			return false
		}
	}
	return true
}

func (s *Server) GetPort() int32 {
	s.PortStatusMutex.RLock()
	defer s.PortStatusMutex.RUnlock()
	for k, v := range s.PortStatus {
		if v == true {
			continue
		} else {
			return k
		}
	}
	return -1
}

func (s *Server) AddConnPort(uid int64, p int32) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	s.ConnPortMap[uid] = p
}

func (s *Server) RemoveConnPort(uid int64) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	delete(s.ConnPortMap, uid)
}

func (s *Server) GetPortByConn(uid int64) int32 {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	return s.ConnPortMap[uid]
}

func (s *Server) GetCurrentCounter() int64 {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	return s.Counter
}

func (s *Server) SendSingle(port int32, count int64) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	serverInstance.ProcessWorker.WorkerStatus[port].Single <- count
}

// AddClientInfo 添加客户端信息
func (s *Server) AddClientInfo(uid int64, username string, port int32) {
	s.ClientInfoMutex.Lock()
	defer s.ClientInfoMutex.Unlock()
	s.ClientInfoMap[uid] = &ClientInfo{
		UID:      uid,
		Username: username,
		Port:     port,
	}
}

// RemoveClientInfo 删除客户端信息
func (s *Server) RemoveClientInfo(uid int64) {
	s.ClientInfoMutex.Lock()
	defer s.ClientInfoMutex.Unlock()
	delete(s.ClientInfoMap, uid)
}

// GetClientInfo 获取客户端信息
func (s *Server) GetClientInfo(uid int64) *ClientInfo {
	s.ClientInfoMutex.RLock()
	defer s.ClientInfoMutex.RUnlock()
	return s.ClientInfoMap[uid]
}
