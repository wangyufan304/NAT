package main

import (
	"net"
	"sync"
)

// Worker 真正干活的工人，数量和TCP MAX有关
type Worker struct {
	// ClientConn 客户端连接
	ClientConn *net.TCPConn
	// ServerListener 服务端监听端口
	ServerListener *net.TCPListener
	// Port 服务端对应端口
	Port int32
	// 该连接对应的user连接池
	TheUserConnPool *userConnPool
	// CurrentTransmitBytes 目前转发了多少个字节
	CurrentTransmitBytes int64
}

type Workers struct {
	Mutex        sync.RWMutex
	WorkerStatus map[int32]*Worker
}

// NewWorkers 新建workers
func NewWorkers() *Workers {
	return &Workers{
		Mutex:        sync.RWMutex{},
		WorkerStatus: make(map[int32]*Worker),
	}
}
func (workers *Workers) Add(port int32, w *Worker) {
	workers.Mutex.Lock()
	defer workers.Mutex.Unlock()
	serverInstance.Counter++
	workers.WorkerStatus[port] = w
	serverInstance.PortStatus[port] = true
}
func (workers *Workers) Remove(port int32) {
	workers.Mutex.Lock()
	defer workers.Mutex.Unlock()
	log.Infoln(workers.WorkerStatus)
	if workers.WorkerStatus[port].ServerListener != nil {
		workers.WorkerStatus[port].ServerListener.Close()
	}
	delete(workers.WorkerStatus, port)
	serverInstance.PortStatus[port] = false
}
func (workers *Workers) Get(port int32) *Worker {
	workers.Mutex.RLock()
	defer workers.Mutex.RUnlock()
	return workers.WorkerStatus[port]
}

func NewWorker(l *net.TCPListener, c *net.TCPConn, port int32) *Worker {
	return &Worker{
		ClientConn:           c,
		ServerListener:       l,
		Port:                 port,
		CurrentTransmitBytes: 0,
		TheUserConnPool:      NewUserConnPool(),
	}
}
