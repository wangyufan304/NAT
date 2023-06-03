package instance

import (
	"github.com/byteYuFan/NAT/ninterfance"
	"net"
)

type Request struct {
	ConnParticularInfo ConnInfo
	Msg                ninterfance.IMessage
}

// NewRequest 新建一个结构体
func NewRequest(ID int64, conn *net.TCPConn, data ninterfance.IMessage) *Request {
	return &Request{
		ConnParticularInfo: ConnInfo{
			Conn: conn,
			ID:   ID,
		},
		Msg: data,
	}
}

// GetConn 得到连接
func (r *Request) GetConn() *net.TCPConn {
	return r.ConnParticularInfo.Conn
}

// GetData 得到请求数据
func (r *Request) GetData() []byte {
	return r.Msg.GetMsgData()
}

// GetID 获取msg ID
func (r *Request) GetID() uint32 {
	return r.Msg.GetMsgID()
}

func (r *Request) GetConnectionUID() int64 {
	return r.ConnParticularInfo.ID
}
