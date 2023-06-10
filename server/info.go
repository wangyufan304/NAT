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
		if err != nil {
			return err
		}
		dbInfo := fmt.Sprintf("%s:%s@tcp(%s)/%s", objectConfig.DB.Username, objectConfig.DB.Password, objectConfig.DB.Host, objectConfig.DB.DBName)
		cui := network.NewControllerUserInfo([]byte(network.KEY), "mysql", dbInfo)
		err = cui.CheckUser(ui)
		if err != nil {
			return err
		}

	}
	return nil
}
