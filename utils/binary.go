package utils

import (
	"github.com/byteYuFan/NAT/ninterfance"
	"io"
	"net"
)

func ObjectToBufferStream(object ninterfance.ObjectTOByteBuffer) ([]byte, error) {
	return object.ToByteBuffer()
}
func BufferStreamToObject(object ninterfance.BufferToObjectInterface, data []byte) error {
	return object.BufferToObject(data)
}
func ReadHeadData(conn *net.TCPConn) ([]byte, error) {
	headData := make([]byte, 8)
	if _, err := io.ReadFull(conn, headData); err != nil {
		return nil, err
	}
	return headData, nil
}

func GetRealData(dataLen uint32, conn *net.TCPConn) ([]byte, error) {
	data := make([]byte, dataLen)
	if _, err := io.ReadFull(conn, data); err != nil {
		return nil, err
	}
	return data, nil
}
