package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "Server-NAT [OPTIONS] COMMAND",
	Short: "Intranet penetration",
	Long:  "GO language-based Intranet penetration tool that supports multiple connections",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var object *objectConfigData

func initCobra() {
	object = &objectConfigData{
		// 初始化对象的字段
	}

	// 将命令行参数与对象的字段绑定
	rootCmd.Flags().StringVarP(&object.Name, "name", "n", "", "Server name")
	rootCmd.Flags().StringVarP(&object.ControllerAddr, "controller-addr", "c", "", "Server controller address")
	rootCmd.Flags().StringVarP(&object.TunnelAddr, "tunnel-addr", "t", "", "Server tunnel address")
	rootCmd.Flags().IntSliceVarP(&object.ExposePort, "expose-port", "p", nil, "Server exposed ports")
	rootCmd.Flags().Int32VarP(&object.TaskQueueNum, "task-queue-num", "q", 0, "Task queue number")
	rootCmd.Flags().Int32VarP(&object.TaskQueueBufferSize, "task-queue-buffer-size", "b", 0, "Task queue buffer size")
	rootCmd.Flags().Int32VarP(&object.MaxTCPConnNum, "max-tcp-conn-num", "m", 0, "Maximum TCP connection number")
	rootCmd.Flags().Int32VarP(&object.MaxConnNum, "max-conn-num", "x", 0, "Maximum connection number")
	rootCmd.Flags().StringVarP(&object.LogFilename, "log-name", "l", "", "The name of the log.")
	rootCmd.Flags().StringVarP(&object.StartAuth, "start-auth", "a", "true", "This is the method that whether the server start the auth.")

	// 打印绑定后的对象

	// 将参数赋值给目标配置对象

	// 添加其他字段...
}

func exchange() {
	if object.Name != "" {
		objectConfig.Name = object.Name
	}
	if object.ControllerAddr != "" {
		objectConfig.ControllerAddr = object.ControllerAddr
	}
	if object.TunnelAddr != "" {
		objectConfig.TunnelAddr = object.TunnelAddr
	}
	if object.LogFilename != "" {
		objectConfig.LogFilename = object.LogFilename
	}
	if object.ExposePort != nil {
		objectConfig.ExposePort = object.ExposePort
	}
	if object.TaskQueueNum != 0 {
		objectConfig.TaskQueueNum = object.TaskQueueNum
	}
	if object.TaskQueueBufferSize != 0 {
		objectConfig.TaskQueueBufferSize = object.TaskQueueBufferSize
	}
	if object.MaxTCPConnNum != 0 {
		objectConfig.MaxTCPConnNum = object.MaxTCPConnNum
	}
	if object.MaxConnNum != 0 {
		objectConfig.MaxConnNum = object.MaxConnNum
	}
	if object.StartAuth != "true" {
		objectConfig.StartAuth = object.StartAuth
	}

}
