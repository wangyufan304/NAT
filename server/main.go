package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

// 这个版本只实现最基本的NAT穿透，即就是最简单的转发
// 流程大概如下
var log = logrus.New()

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		// 只打印帮助信息，不执行命令
		Execute()
	} else {
		Execute()
		art()
		exchange()
		printServerRelationInformation()
		go ListenTaskQueue()
		// go acceptClientRequest()
		go cleanExpireConnPool()
		createControllerChannel()
	}

}

func init() {
	objectConfig = new(objectConfigData)
	initConfig()
	initCobra()
	initLog()
	initServer()
	initUserConnPool()
}
