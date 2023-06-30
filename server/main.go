package main

import (
	"github.com/byteYuFan/NAT/instance"
	"github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
	"os"
)

// 这个版本只实现最基本的NAT穿透，即就是最简单的转发
// 流程大概如下
var log = logrus.New()
var myLogger *instance.MyLogger

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		// 只打印帮助信息，不执行命令
		Execute()
	} else {
		Execute()
		art()
		exchange()
		initLogger()
		printServerRelationInformation()
		go func() {
			log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
		}()
		go createControllerChannel()
		go ListenTaskQueue()
		go acceptClientRequest()
		select {}
	}

}

func init() {
	objectConfig = new(objectConfigData)
	initConfig()
	initCobra()
	initLog()
	initServer()
}
