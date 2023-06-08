package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

func initLog() {
	file, err := os.OpenFile(objectConfig.LogFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		// 设置日志输出的文件
		log.SetOutput(file)
	} else {
		log.Fatal("Failed to open log file:", err)
	}
	// 设置日志级别
	log.SetLevel(logrus.InfoLevel)
}
