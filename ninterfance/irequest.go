package ninterfance

import "net"

// IRequest 请求数据
type IRequest interface {
	GetConn() *net.TCPConn
	GetData() []byte
	GetID() uint32
}
