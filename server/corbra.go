package main

import "github.com/spf13/cobra"

func initCobra() {
	object := new(objectConfigData)
	rootCmd.Flags().StringVarP(&object.Name, "name", "n", "", "Server name")
	rootCmd.Flags().StringVarP(&object.ControllerAddr, "controller-addr", "c", "", "Server controller address")
	rootCmd.Flags().StringVarP(&object.TunnelAddr, "tunnel-addr", "t", "", "Server tunnel address")
	rootCmd.Flags().IntSliceVarP(&object.ExposePort, "expose-port", "p", nil, "Server exposed ports")
	rootCmd.Flags().Int32VarP(&object.TaskQueueNum, "task-queue-num", "q", 0, "Task queue number")
	rootCmd.Flags().Int32VarP(&object.TaskQueueBufferSize, "task-queue-buffer-size", "b", 0, "Task queue buffer size")
	rootCmd.Flags().Int32VarP(&object.MaxTCPConnNum, "max-tcp-conn-num", "m", 0, "Maximum TCP connection number")
	rootCmd.Flags().Int32VarP(&object.MaxConnNum, "max-conn-num", "x", 0, "Maximum connection number")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if object.Name != "" {
			objectConfig.Name = object.Name
		}
		if object.ControllerAddr != "" {
			objectConfig.ControllerAddr = object.ControllerAddr
		}
		if object.TunnelAddr != "" {
			objectConfig.TunnelAddr = object.TunnelAddr
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
		return nil
	}

	// 添加其他字段...
}
