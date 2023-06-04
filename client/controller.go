package main

import (
	"github.com/byteYuFan/NAT/instance"
	"github.com/byteYuFan/NAT/ninterfance"
	"io"
	"net"
)

// ControllerService  这是控制层连接的一个实体
// 控制层的主要行为为 向客户端发送信息
type ControllerService struct {
	// Conn 控制端的连接socket
	Conn *net.TCPConn
}

// NewControllerServiceInstance 新建一个控制层实例对象
func NewControllerServiceInstance(conn *net.TCPConn) *ControllerService {
	return &ControllerService{
		Conn: conn,
	}
}

// SendDataToClient 向客户端发送消息，此处应该指定协议
func (cs *ControllerService) SendDataToClient(dataType uint32, msg []byte) (int, error) {
	msgInstance := instance.NewMsgPackage(dataType, msg)
	pkgInstance := instance.NewDataPackage()
	dataStream, err := pkgInstance.Pack(msgInstance)
	if err != nil {
		return 0, err
	}
	count, err := cs.Conn.Write(dataStream)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (cs *ControllerService) ReadHeadDataFromClient() (ninterfance.IMessage, error) {
	// 根据我们协议规定的内容每次数据来的时候，先读取头部的8个字节的数据
	headData := make([]byte, 8)
	// 从Conn中读取数据到headData中去
	if _, err := io.ReadFull(cs.Conn, headData); err != nil {
		return nil, err
	}
	// 先创建一个解包的实例
	dp := instance.NewDataPackage()
	// 解封装这个包
	return dp.Unpack(headData)
}
func (cs *ControllerService) ReadRealDataFromClient(msg ninterfance.IMessage) (ninterfance.IMessage, error) {
	if msg.GetMsgDataLen() < 0 {
		return msg, nil
	}
	// 新建一个data，长度为msg头部的长度
	data := make([]byte, msg.GetMsgDataLen())
	if _, err := io.ReadFull(cs.Conn, data); err != nil {
		return nil, err
	}
	msg.SetMsgData(data)
	return msg, nil
}
