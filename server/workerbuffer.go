package main

import "net"

type Request struct {
	Conn     *net.TCPConn
	Username string
	ID       int64
}
