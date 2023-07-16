package main

import (
	"errors"
	"fmt"
	"github.com/byteYuFan/NAT/instance"
	"github.com/byteYuFan/NAT/network"
	"net"
)

func authUser(conn *net.TCPConn) (*network.UserInfo, error) {
	nsi := instance.NewSendAndReceiveInstance(conn)
	msg, err := nsi.ReadHeadDataFromClient()
	if err != nil {
		return nil, err
	}
	msg, err = nsi.ReadRealDataFromClient(msg)
	if err != nil {
		return nil, err
	}
	ui := new(network.UserInfo)
	if msg.GetMsgID() == network.USER_REQUEST_AUTH {
		// 获取其真实数据
		err := ui.FromBytes(msg.GetMsgData())
		if err != nil {
			return nil, err
		}
		dbInfo := fmt.Sprintf("%s:%s@tcp(%s)/%s", objectConfig.DB.Username, objectConfig.DB.Password, objectConfig.DB.Host, objectConfig.DB.DBName)
		fmt.Println(dbInfo)
		fmt.Println(ui)
		cui := network.NewControllerUserInfo([]byte(network.KEY), "mysql", dbInfo)
		if cui == nil {
			return nil, errors.New("初始化错误")
		}
		err = cui.CheckUser(ui)
		if err != nil {
			return nil, err
		}

	}
	return ui, nil
}
