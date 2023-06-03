package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		fmt.Println("Failed to connect:", err)
		return
	}
	fmt.Println("[连接成功]")
	defer conn.Close()

	// 发送黏包数据
	message := "Hello!"
	msgOne := MessageProtocol{
		Len:  uint32(len(message)),
		Data: []byte(message),
	}

	byteBuffer, err := pack(msgOne)
	fmt.Println("[byteBuffer]", byteBuffer)
	if err != nil {
		fmt.Println("[pack]", err)
	}
	writen, err := conn.Write(byteBuffer)
	if err != nil {
		fmt.Println("[write ]", err.Error())
		return
	}
	fmt.Println("[write successfully]", writen)
	message = "How are you?"
	msgOne = MessageProtocol{
		Len:  uint32(len(message)),
		Data: []byte(message),
	}
	byteBuffer, _ = pack(msgOne)
	fmt.Println("[byteBuffer]", byteBuffer)
	writen, _ = conn.Write(byteBuffer)
	fmt.Println("[write successfully]", writen)
}
