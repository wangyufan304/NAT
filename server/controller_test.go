package main

import (
	"fmt"
	"github.com/byteYuFan/NAT/network"
	"testing"
	"time"
)

func TestControllerService_ReadHeadDataFromClient(t *testing.T) {
	conn, _ := network.CreateTCPListener("127.0.0.1:8080")
	go func() {
		for {
			tcp, _ := conn.AcceptTCP()
			nsi := NewControllerServiceInstance(tcp)
			for {
				time.Sleep(time.Second)
				client, err := nsi.SendDataToClient(1, []byte("hello,how are you"))
				if err != nil {
					fmt.Println("[Send]", err)
					return
				}
				fmt.Println("[Send]", client)
			}
		}
	}()
	go func() {

		clientConn, _ := network.CreateTCPConn("127.0.0.1:8080")
		nsi := NewControllerServiceInstance(clientConn)
		for {
			msg, err := nsi.ReadHeadDataFromClient()
			if err != nil {
				fmt.Println("[readHead]", msg)
				continue
			}
			msg, err = nsi.ReadRealDataFromClient(msg)
			if err != nil {
				fmt.Println("[readReal]", msg)
				continue
			}
			fmt.Println("Read", string(msg.GetMsgData()))
		}
	}()
	select {}
}
