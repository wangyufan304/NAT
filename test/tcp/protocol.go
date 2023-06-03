package main

import (
	"bytes"
	"encoding/binary"
)

type MessageProtocol struct {
	// Len 数据长度，注意此处必须为无符号类型
	Len uint32
	// Data 真实数据
	Data []byte
}

func pack(msg MessageProtocol) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.Len); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.Data); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

func unPack(data []byte) (*MessageProtocol, error) {
	dataBuffer := bytes.NewBuffer(data)
	msg := new(MessageProtocol)
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.Len); err != nil {
		return nil, err
	}
	return msg, nil
}
