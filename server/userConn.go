package main

import (
	"net"
	"strconv"
	"sync"
	"time"
)

// UserConnInfo 用户连接信息,此处保存的是用户访问公网web所对应的那个接口
type UserConnInfo struct {
	// visit 用户访问web服务的时间
	visit time.Time
	// conn tcp连接句柄
	conn *net.TCPConn
}

// userConnPool 用户连接池
type userConnPool struct {
	// UserConnectionMap 连接池map，存放着用户的连接信息 key-时间戳 val userConnInfo
	UserConnectionMap map[string]*UserConnInfo
	// Mutex 读写锁，用来保证map的并发安全问题
	Mutex sync.RWMutex
}

var userConnPoolInstance *userConnPool

// NewUserConnPool 新建一个连接池对象
func NewUserConnPool() *userConnPool {
	return &userConnPool{
		UserConnectionMap: make(map[string]*UserConnInfo),
		Mutex:             sync.RWMutex{},
	}
}

// AddConnInfo 向连接池添加用户的信息
func (ucp *userConnPool) AddConnInfo(conn *net.TCPConn) {
	// 加写锁保护并发的安全性
	ucp.Mutex.Lock()
	defer ucp.Mutex.Unlock()
	nowTime := time.Now()
	uci := &UserConnInfo{
		visit: nowTime,
		conn:  conn,
	}
	ucp.UserConnectionMap[strconv.FormatInt(nowTime.UnixNano(), 10)] = uci
}
func initUserConnPool() {
	userConnPoolInstance = NewUserConnPool()
}
