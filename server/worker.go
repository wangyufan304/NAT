package main

import "net"

// WorkerQueue 工作队列
type WorkerQueue struct {
	// Worker 工具人
	Worker chan *net.TCPConn
	// Port 与该队列绑定的服务端Port
	Port int
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
