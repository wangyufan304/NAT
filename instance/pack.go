package instance

import (
	"bytes"
	"encoding/binary"
	"github.com/byteYuFan/NAT/ninterfance"
)

// DataPackage 封包解包的结构体
type DataPackage struct {
}

// NewDataPackage 创建一个封包拆包的实例
func NewDataPackage() *DataPackage {
	return &DataPackage{}
}

// GetHeadLen 获取包头的长度 根据我们的协议定义直接返回8就可以了
func (dp *DataPackage) GetHeadLen() uint32 {
	return uint32(8)
}

// Pack 将 ninterfance.IMessage 类型的结构封装为字节流的形式
// 字节流形式 [ 数据长度 + ID + 真实数据 ]
func (dp *DataPackage) Pack(msg ninterfance.IMessage) ([]byte, error) {
	// 创建一个字节流的缓存，将msg的信息一步一步的填充到里面去
	dataBuff := bytes.NewBuffer([]byte{})
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetMsgDataLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetMsgData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

func (dp *DataPackage) Unpack(data []byte) (ninterfance.IMessage, error) {
	// 创建一个从data里面读取的ioReader
	dataBuffer := bytes.NewBuffer(data)
	msg := &Message{}
	if err := binary.Read(dataBuffer, binary.BigEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuffer, binary.BigEndian, &msg.ID); err != nil {
		return nil, err
	}
	return msg, nil
}
