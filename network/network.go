package network

import (
	"io"
	"net"
	"sync"
)

const (
	NewConnection = "NEW_CONNECTION"
)

// 这个包里面放着一些通用的函数

// SwapConnDataEachOther 通讯双方相互交换数据
//func SwapConnDataEachOther(local, remote *net.TCPConn) {
//	go swapConnData(local, remote)
//	go swapConnData(remote, local)
//}

// SwapConnData 这个函数是交换两个连接数据的函数
//func swapConnData(local, remote *net.TCPConn) {
//	// 关闭本地和远程连接通道
//	defer local.Close()
//	defer remote.Close()
//	// 将remote的数据拷贝到local里面
//	_, err := io.Copy(local, remote)
//	if err != nil {
//		return
//	}
//}

// SwapConnDataEachOther 通讯双方相互交换数据
func SwapConnDataEachOther(local, remote *net.TCPConn) int64 {
	wg := sync.WaitGroup{}
	wg.Add(2)
	var count1 int64
	var count2 int64
	go func(count *int64) {
		*count = swapConnData(local, remote, &wg)
	}(&count1)
	go func(count *int64) {
		*count = swapConnData(remote, local, &wg)
	}(&count2)
	wg.Wait()
	return count1 + count2
}

// swapConnData 这个函数是交换两个连接数据的函数
func swapConnData(src, dst *net.TCPConn, wg *sync.WaitGroup) int64 {
	defer wg.Done()

	// 创建一个固定大小的缓冲区用于数据拷贝
	buffer := make([]byte, 1024)
	var count int64
	for {
		n, err := src.Read(buffer)
		count += int64(n)
		if err != nil {
			if err != io.EOF {
			}
			break
		}
		_, err = dst.Write(buffer[:n])
		if err != nil {
			break
		}
	}
	return count
}

func CreateTCPListener(addr string) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}
	return tcpListener, nil
}

// CreateTCPConn 连接指定的TCP
func CreateTCPConn(addr string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	return tcpConn, nil
}
