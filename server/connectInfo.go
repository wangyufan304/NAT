package main

import (
	"net"
	"sync"
)

// serverConnInfo 连接信息模块,该结构存储着所有的连接信息
type serverConnInfo struct {
	// Port 该连接对应的端口号
	Port int32
	// Conn 该链接的tcp socket
	Conn *net.TCPConn
	// Listener 该连接对应的 socket listener
	Listener *net.TCPListener
	// Count
	Count int64
	// Username
	Username string
}

var GlobalInfo *Information

func NewServerConnInfo(port int32, conn *net.TCPConn, listener *net.TCPListener, username string) *serverConnInfo {
	return &serverConnInfo{
		Port:     port,
		Conn:     conn,
		Listener: listener,
		Count:    0,
		Username: username,
	}
}

type Information struct {
	infoMAP map[int64]*serverConnInfo
	sync.RWMutex
}

func (i *Information) Add(uid int64, port int32, conn *net.TCPConn, listener *net.TCPListener, username string) {
	i.RWMutex.Lock()
	defer i.RWMutex.Unlock()
	i.infoMAP[uid] = NewServerConnInfo(port, conn, listener, username)
}

func (i *Information) Remove(uid int64) {
	i.RWMutex.Lock()
	defer i.RWMutex.Unlock()
	delete(i.infoMAP, uid)
}

func (i *Information) GetServerInfo(uid int64) *serverConnInfo {
	i.RWMutex.RLock()
	defer i.RWMutex.RUnlock()
	return i.infoMAP[uid]
}

func init() {
	GlobalInfo = &Information{
		infoMAP: make(map[int64]*serverConnInfo),
		RWMutex: sync.RWMutex{},
	}
}
