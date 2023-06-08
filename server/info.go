package main

import (
	"fmt"
	"github.com/byteYuFan/NAT/instance"
	"github.com/byteYuFan/NAT/network"
	"net"
)

func authUser(conn *net.TCPConn) error {
	nsi := instance.NewSendAndReceiveInstance(conn)
	msg, err := nsi.ReadHeadDataFromClient()
	if err != nil {
		return err
	}
	msg, err = nsi.ReadRealDataFromClient(msg)
	if err != nil {
		return err
	}
	if msg.GetMsgID() == network.USER_REQUEST_AUTH {
		// 获取其真实数据
		ui := new(network.UserInfo)
		err := ui.FromBytes(msg.GetMsgData())
		fmt.Println("[ui]", ui)
		if err != nil {
			return err
		}
		cui := network.NewControllerUserInfo([]byte(network.KEY), "mysql", "root:123456@tcp(pogf.com.cn:3309)/NAT")
		err = cui.CheckUser(ui)
		if err != nil {
			return err
		}

	}
	return nil
}
