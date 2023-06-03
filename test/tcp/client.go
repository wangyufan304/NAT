package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	ln, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("Failed to listen:", err)
		return
	}
	defer ln.Close()

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Failed to accept connection:", err)
		return
	}
	fmt.Println("[conn]", conn.RemoteAddr().String())
	for {
		if conn == nil {
			break
		}
		// 首先从字节流中读取4 字节的内容，为什么读取4 字节当然是我们协议的规定啦
		var headLen = make([]byte, 4)
		if _, err := io.ReadFull(conn, headLen); err != nil {
			fmt.Println("[ReadFull]", err)
			return
		}
		msg, err := unPack(headLen)
		if err != nil {
			fmt.Println("[unPack]", err)
			return
		}
		var data []byte
		if msg.Len > 0 {
			data = make([]byte, msg.Len)
			if _, err := io.ReadFull(conn, data); err != nil {
				fmt.Println("[ReadFull]", err)
				return
			}
		}
		fmt.Println("[READDATA]", string(data), "[BYTE]:", data)
		time.Sleep(time.Second)
	}
	fmt.Println("[Conn]", "客户端关闭")
}
