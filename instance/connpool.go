package instance

import (
	"net"
	"sync"
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
			Port:             []uint32{60001, 60002, 600003, 60004},
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
