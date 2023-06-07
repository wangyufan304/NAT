package main

import (
	"net"
	"sync"
)

// WorkerQueue 工作队列
type WorkerQueue struct {
	// Worker 工具人
	Worker chan *net.TCPConn
	// Port 与该队列绑定的服务端Port
	Port int
}

// Worker 真正干活的工人，数量和TCP MAX有关
type Worker struct {
	// ClientConn 客户端连接
	ClientConn *net.TCPConn
	// ServerListener 服务端监听端口
	ServerListener *net.TCPListener
	// Port 服务端对应端口
	Port int32
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
	workers.WorkerStatus[port] = w
	serverInstance.PortStatus[port] = true
}
func (workers *Workers) Remove(port int32) {
	workers.Mutex.Lock()
	defer workers.Mutex.Unlock()
	workers.WorkerStatus[port].ServerListener.Close()
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
		ClientConn:     c,
		ServerListener: l,
		Port:           port,
	}
}

// NewWorkerQueue 新建一个工作队列
func NewWorkerQueue(bufferSize int32, port int) *WorkerQueue {
	return &WorkerQueue{
		Worker: make(chan *net.TCPConn, bufferSize),
		Port:   port,
	}
}

func (wq *WorkerQueue) GetPort() int {
	return wq.Port
}
