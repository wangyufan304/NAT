package instance

import (
	"net"
	"sync"
	"sync/atomic"
)

type ConnPool struct {
	MaxTCPConn       uint32
	TaskQueue        []chan *Request
	BufferSize       uint32
	Port             []uint32
	PortAndTaskQueue map[uint32]uint32
	ConnInfo         map[*net.TCPConn]uint32
	Counter          int64
	CurrentConnNum   int32
	sync.RWMutex
}

var (
	ConnPoolInstance *ConnPool
	ConnPoolInitOnce sync.Once
)

func NewConnPool() *ConnPool {
	ConnPoolInitOnce.Do(func() {
		ConnPoolInstance = &ConnPool{
			MaxTCPConn:       4,
			TaskQueue:        make([]chan *Request, 4),
			BufferSize:       8,
			Port:             []uint32{60000, 60001, 600002, 60003},
			PortAndTaskQueue: make(map[uint32]uint32),
			ConnInfo:         make(map[*net.TCPConn]uint32),
			Counter:          0,
		}
	})
	for i := 0; i < int(ConnPoolInstance.MaxTCPConn); i++ {
		// TODO
	}
	return ConnPoolInstance
}

func (cp *ConnPool) GetMaxTCPConn() uint32 {
	return cp.MaxTCPConn
}

func (cp *ConnPool) AddTCPConn(conn *net.TCPConn) int64 {
	// 当接收到新地连接请求时
	cp.Lock()
	defer cp.Unlock()
	req := NewRequest(cp.Counter, conn, nil)
	atomic.AddInt64(&cp.Counter, 1)
	// 采用轮询的方式为客户端分配uid 采用与最大连接数取余作为负载均衡
	uid := req.GetConnectionUID() % int64(cp.GetMaxTCPConn())
	cp.TaskQueue[uid] <- req
	cp.ConnInfo[conn] = cp.Port[uid]
	return uid
}

func (cp *ConnPool) RemoveConn(uid int64, conn *net.TCPConn) *Request {
	cp.Lock()
	defer cp.Unlock()
	delete(cp.ConnInfo, conn)
	cp.Counter--
	return <-cp.TaskQueue[uid]
}
