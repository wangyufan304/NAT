package network

import (
	"bytes"
	"encoding/binary"
)

// ClientConnInfo 客户端连接信息
type ClientConnInfo struct {
	UID  int64
	Port int32
}

// NewClientConnInstance 新建一个实体
func NewClientConnInstance(id int64, port int32) *ClientConnInfo {
	return &ClientConnInfo{
		UID:  id,
		Port: port,
	}
}

// ToBytes 将 ClientConnInfo 结构体转换为字节流
func (info *ClientConnInfo) ToBytes() ([]byte, error) {
	buf := new(bytes.Buffer)

	// 使用 binary.Write 将字段逐个写入字节流
	if err := binary.Write(buf, binary.BigEndian, info.UID); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, info.Port); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// FromBytes 从字节流中恢复 ClientConnInfo 结构体
func (info *ClientConnInfo) FromBytes(data []byte) error {
	buf := bytes.NewReader(data)
	// 使用 binary.Read 从字节流中读取字段值
	if err := binary.Read(buf, binary.BigEndian, &info.UID); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &info.Port); err != nil {
		return err
	}
	return nil
}
