package instance

import (
	"io"
	"log"
	"os"
)

// 自己封装一个log库

// MyLogger 自己简单封装的一个日志结构，该日志有三个column
// 普通信息日志 成功信息日志 错误信息日志
type MyLogger struct {
	InfoLogger    *log.Logger
	SuccessLogger *log.Logger
	ErrorLogger   *log.Logger
}

// NewMyLogger 实体
func NewMyLogger(out io.Writer, fileOutput bool, filename string) *MyLogger {
	logger := &MyLogger{}

	// 设置输出目标
	if fileOutput {
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("无法打开日志文件：%v", err)
		}
		out = file

		// 创建错误日志文件
		errorFile, err := os.OpenFile("errorlog.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("无法打开错误日志文件：%v", err)
		}

		// 设置错误日志记录器的输出目标
		logger.ErrorLogger = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime)
	} else {
		// 设置终端为日志记录器的输出目标
		out = os.Stdout
	}

	// 配置日志格式和输出目标
	logger.InfoLogger = log.New(out, "INFO: ", log.Ldate|log.Ltime)
	logger.SuccessLogger = log.New(out, "SUCCESS: ", log.Ldate|log.Ltime)

	return logger
}

func (l *MyLogger) Info(message string) {
	l.InfoLogger.Println(message)
}

func (l *MyLogger) Success(message string) {
	l.SuccessLogger.Println(message)
}

func (l *MyLogger) Error(message string) {
	l.ErrorLogger.Println(message)
}
