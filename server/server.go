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
	WorkerBuffer chan *net.TCPConn
	// 实际处理工作的数据结构
	ProcessWorker *Workers
	// 端口使用情况
	PortStatus map[int32]bool
}

func initServer() {
	serverInstance = &Server{
		Mutex:          sync.RWMutex{},
		Counter:        0,
		MaxTCPConnSize: objectConfig.MaxTCPConnNum,
		MaxConnSize:    objectConfig.MaxConnNum,
		ExposePort:     objectConfig.ExposePort,
		ProcessingMap:  make(map[string]*net.TCPConn),
		WorkerBuffer:   make(chan *net.TCPConn, objectConfig.MaxConnNum),
		ProcessWorker:  NewWorkers(),
		PortStatus:     make(map[int32]bool),
	}

	// 初始化端口状态
	for i := 0; i < int(serverInstance.MaxTCPConnSize); i++ {
		serverInstance.PortStatus[int32(serverInstance.ExposePort[i])] = false
	}
}

func (s *Server) PortIsFull() bool {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	for _, v := range s.PortStatus {
		if v == false {
			return false
		}
	}
	return true
}

func (s *Server) GetPort() int32 {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	for k, v := range s.PortStatus {
		if v == true {
			continue
		} else {
			return k
		}
	}
	return -1
}
