package main

import (
	"github.com/byteYuFan/NAT/utils"
)

// ObjectConfigData 全局配置的对象,里面存储着服务端所有的配置信息
type objectConfigData struct {
	// ServerName 服务端名称
	Name string
	// ControllerAddr 服务器控制端地址
	ControllerAddr string
	// TunnelAddr 隧道地址交换数据隧道
	TunnelAddr string
	// ExposePort 服务端向外暴露的端口号
	ExposePort []int
	// TaskQueueNum 任务队列的数量
	TaskQueueNum int32
	// TaskQueueBufferSize 缓冲区最大的数量
	TaskQueueBufferSize int32
	// MaxTCPConnNum  一次性最大处理的并发连接数量，等同于任务队列的大小和服务端暴露的端口号数量
	MaxTCPConnNum int32
	// MaxConnNum 整个系统所能接收到的并发数量 为工作队列大小和工作队列缓冲区之积
	MaxConnNum int32
	// 	LogFilename string 日志文件名称
	LogFilename string
	// StartAuth 是否开启认证功能
	StartAuth string
	// DB 如果开启认证功能就得从配置文件中读取相关的配置信息
	DB DataBase
}

// DataBase 数据库相关信息
type DataBase struct {
	Username string
	Password string
	Host     string
	DBName   string
}

var objectConfig *objectConfigData

func initConfig() {
	// 读取配置文件内容
	config := utils.ParseFile("server.yml")
	viper := config.ViperInstance
	viper.SetDefault("Server.Name", "Server-NAT")
	viper.SetDefault("Server.ControllerAddr", "0.0.0.0:8007")
	viper.SetDefault("Server.TunnelAddr", "0.0.0.0:8008")
	viper.SetDefault("Server.VisitPort", []uint16{60000, 60001, 60002, 60003})
	viper.SetDefault("Server.TaskQueueNum", 4)
	viper.SetDefault("Server.TaskQueueBuff", 32)
	viper.SetDefault("Server.MaxTCPConnNum", 4)
	viper.SetDefault("Server.MaxConnNum", 128)
	viper.SetDefault("Server.LogFilename", "server.log")
	viper.SetDefault("Server.StartAuth", true)
	// 读取配置值并存入 objectConfig
	objectConfig.Name = viper.GetString("Server.Name")
	objectConfig.ControllerAddr = viper.GetString("Server.ControllerAddr")
	objectConfig.TunnelAddr = viper.GetString("Server.TunnelAddr")
	objectConfig.ExposePort = viper.GetIntSlice("Server.VisitPort")
	objectConfig.TaskQueueNum = viper.GetInt32("Server.TaskQueueNum")
	objectConfig.MaxTCPConnNum = viper.GetInt32("Server.MaxTCPConnNum")
	objectConfig.TaskQueueBufferSize = viper.GetInt32("Server.TaskQueueBuff")
	objectConfig.MaxConnNum = viper.GetInt32("Server.MaxConnNum")
	objectConfig.LogFilename = viper.GetString("Server.LogFilename")
	objectConfig.StartAuth = viper.GetString("Server.StartAuth")
	objectConfig.DB.Username = viper.GetString("Database.Username")
	objectConfig.DB.Password = viper.GetString("Database.Password")
	objectConfig.DB.Host = viper.GetString("Database.Host")
	objectConfig.DB.DBName = viper.GetString("Database.DBName")
}
