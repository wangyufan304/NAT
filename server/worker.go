package main

import (
	"net"
	"sync"
)

// Worker 真正干活的工人，数量和TCP MAX有关
type Worker struct {
	// 对应客户端的ID
	ID int64
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
	// Single 信号带有8个缓冲区的buffer
	Single chan int64
	mutex  sync.RWMutex
}

type Workers struct {
	Mutex      sync.RWMutex
	WorkerInfo map[int64]*Worker
}

// NewWorkers 新建workers
func NewWorkers() *Workers {
	return &Workers{
		Mutex:      sync.RWMutex{},
		WorkerInfo: make(map[int64]*Worker),
	}
}
func (workers *Workers) Add(uid int64, w *Worker) {
	workers.Mutex.Lock()
	defer workers.Mutex.Unlock()
	workers.WorkerInfo[uid] = w
}

func (workers *Workers) Remove(uid int64) {
	workers.Mutex.Lock()
	defer workers.Mutex.Unlock()
	delete(workers.WorkerInfo, uid)
}

func (workers *Workers) Get(uid int64) *Worker {
	workers.Mutex.RLock()
	defer workers.Mutex.RUnlock()
	return workers.WorkerInfo[uid]
}

func NewWorker(l *net.TCPListener, c *net.TCPConn, port int32, UID int64) *Worker {
	return &Worker{
		ID:                   UID,
		ClientConn:           c,
		ServerListener:       l,
		Port:                 port,
		CurrentTransmitBytes: 0,
		TheUserConnPool:      NewUserConnPool(),
		Single:               make(chan int64, 8),
		mutex:                sync.RWMutex{},
	}
}
func (w *Worker) GetTheUserConnPool() *userConnPool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.TheUserConnPool
}

func (w *Worker) GetCounter() int64 {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.CurrentTransmitBytes
}

func (w *Worker) AddCount(num int64) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.CurrentTransmitBytes += num
}
