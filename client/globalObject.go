package main

import "github.com/byteYuFan/NAT/utils"

type objectConfigData struct {
	// Name 客户端名称
	Name string
	// LocalServerAddr 本地服务端地址
	LocalServerAddr string
	// TunnelServerAddr 隧道地址用于交换数据
	TunnelServerAddr string
	// PublicServerAddr 公网服务器地址
	PublicServerAddr string
	// ControllerAddr 服务器控制端地址
	ControllerAddr string
	// UserName 登录用户名
	UserName string
	// Password 密码
	Password string
}

var objectConfig *objectConfigData

func initConfig() {
	objectConfig = new(objectConfigData)
	config := utils.ParseFile("client.yml")
	viper := config.ViperInstance
	viper.SetDefault("Client.Name", "Client-NAT")
	viper.SetDefault("Client.PublicServerAddr", "pogf.com.cn")
	viper.SetDefault("Client.TunnelServerAddr", "pogf.com.cn:8008")
	viper.SetDefault("Client.ControllerAddr", "pogf.com.cn:8080")
	viper.SetDefault("Client.LocalServerAddr", "127.0.0.1:8080")
	viper.SetDefault("Auth.Username", "")
	viper.SetDefault("Auth.password", "")

	objectConfig.Name = viper.GetString("Client.Name")
	objectConfig.PublicServerAddr = viper.GetString("Client.PublicServerAddr")
	objectConfig.ControllerAddr = viper.GetString("Client.ControllerAddr")
	objectConfig.LocalServerAddr = viper.GetString("Client.LocalServerAddr")
	objectConfig.TunnelServerAddr = viper.GetString("Client.TunnelServerAddr")
	objectConfig.UserName = viper.GetString("Auth.Username")
	objectConfig.Password = viper.GetString("Auth.Password")
}
