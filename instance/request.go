package instance

import "net"

// Request 内网客户端请求的数据
type Request struct {
	// Conn 自己的连接信息
	Conn *net.TCPConn
}
