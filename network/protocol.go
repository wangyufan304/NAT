package network

import (
	"bytes"
	"encoding/binary"
)

// 该文件里面存储着所有协议的定义

// ControllerInfo 控制信息存放着客户端的uid和对应的port
// 该端口的意思就是用户通过该端口可以访问到内网的web等服务
type ControllerInfo struct {
	ID             uint64
	Port           uint32
	CurrentConnNum uint32
}

// KeepAlive 心跳结构体
type KeepAlive struct {
	ID  uint64
	Msg []byte
}

func NewKeepAlive(id uint64, msg string) *KeepAlive {
	return &KeepAlive{
		ID:  id,
		Msg: []byte(msg),
	}
}
func (ci *ControllerInfo) ToByteBuffer() ([]byte, error) {
	buffer := new(bytes.Buffer)
	if err := binary.Write(buffer, binary.LittleEndian, ci); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
func (ci *ControllerInfo) BufferToObject(data []byte) error {
	buffer := bytes.NewBuffer(data)
	if err := binary.Read(buffer, binary.LittleEndian, ci); err != nil {
		return err
	}
	return nil
}
func (ka *KeepAlive) ToByteBuffer() ([]byte, error) {
	buffer := new(bytes.Buffer)
	if err := binary.Write(buffer, binary.LittleEndian, ka.ID); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.LittleEndian, ka.Msg); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
func (ka *KeepAlive) BufferToObject(data []byte) error {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.LittleEndian, &ka.ID); err != nil {
		return err
	}
	if err := binary.Read(buffer, binary.LittleEndian, ka.Msg); err != nil {
		return err
	}
	return nil
}
