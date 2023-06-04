package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// 这个版本只实现最基本的NAT穿透，即就是最简单的转发
// 流程大概如下

var rootCmd = &cobra.Command{
	Use:   "Server-NAT [OPTIONS] COMMAND",
	Short: "Intranet penetration",
	Long:  "GO language-based Intranet penetration tool that supports multiple connections",
	Run: func(cmd *cobra.Command, args []string) {
		// 运行命令的处理逻辑
	},
}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		// 只打印帮助信息，不执行命令
		rootCmd.SetArgs(os.Args[1:])
		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		go createControllerChannel()
		go ListenTaskQueue()
		go acceptClientRequest()
		go cleanExpireConnPool()
		select {}
	}

}

func init() {
	initConfig()
	initCobra()
	initServer()
	initUserConnPool()
}
