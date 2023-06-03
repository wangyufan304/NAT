package instance

import "net"

type ConnInfo struct {
	Conn *net.TCPConn
	ID   int64
}
