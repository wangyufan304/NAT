package instance

import (
	"github.com/byteYuFan/NAT/ninterfance"
	"io"
	"net"
)

// SendAndReceiveInstance 实体
type SendAndReceiveInstance struct {
	Conn *net.TCPConn
}

// NewSendAndReceiveInstance 新建一个控制层实例对象
func NewSendAndReceiveInstance(conn *net.TCPConn) *SendAndReceiveInstance {
	return &SendAndReceiveInstance{
		Conn: conn,
	}
}

// SendDataToClient 向客户端发送消息，此处应该指定协议
func (csi *SendAndReceiveInstance) SendDataToClient(dataType uint32, msg []byte) (int, error) {
	msgInstance := NewMsgPackage(dataType, msg)
	pkgInstance := NewDataPackage()
	dataStream, err := pkgInstance.Pack(msgInstance)
	if err != nil {
		return 0, err
	}
	count, err := csi.Conn.Write(dataStream)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (csi *SendAndReceiveInstance) ReadHeadDataFromClient() (ninterfance.IMessage, error) {
	// 根据我们协议规定的内容每次数据来的时候，先读取头部的8个字节的数据
	headData := make([]byte, 8)
	// 从Conn中读取数据到headData中去
	if _, err := io.ReadFull(csi.Conn, headData); err != nil {
		return nil, err
	}
	// 先创建一个解包的实例
	dp := NewDataPackage()
	// 解封装这个包
	return dp.Unpack(headData)
}
func (csi *SendAndReceiveInstance) ReadRealDataFromClient(msg ninterfance.IMessage) (ninterfance.IMessage, error) {
	if msg.GetMsgDataLen() < 0 {
		return msg, nil
	}
	// 新建一个data，长度为msg头部的长度
	data := make([]byte, msg.GetMsgDataLen())
	if _, err := io.ReadFull(csi.Conn, data); err != nil {
		return nil, err
	}
	msg.SetMsgData(data)
	return msg, nil
}
