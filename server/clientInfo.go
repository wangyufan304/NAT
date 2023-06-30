package main

// ClientInfo 客户端连接信息
type ClientInfo struct {
	// UID 客户端唯一ID
	UID int64
	// Username
	Username string
	// Port 服务器给该客户端分配的端口号
	Port int32
}
