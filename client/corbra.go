package main

import "github.com/spf13/cobra"

func initCobra() {
	rootCmd.Flags().StringVarP(&objectConfig.Name, "name", "n", "Client-NAT", "Client name")
	rootCmd.Flags().StringVarP(&objectConfig.LocalServerAddr, "local-server-addr", "l", "0.0.0.0:8080", "The address of the local web server program")
	rootCmd.Flags().StringVarP(&objectConfig.TunnelServerAddr, "tunnel-server-addr", "t", "pogf.com.cn:8008", "The address of the tunnel server used to connect the local and public networks")
	rootCmd.Flags().StringVarP(&objectConfig.ControllerAddr, "controller-addr", "c", "pogf.com.cn:8080", "The address of the controller channel used to send controller messages to the client")
	rootCmd.Flags().StringVarP(&objectConfig.PublicServerAddr, "public-server-addr", "s", "pogf.com.cn", "The address of the public server used for accessing the inner web server")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if name := cmd.Flag("name").Value.String(); name != "" {
			objectConfig.Name = name
		}
		if localServerAddr := cmd.Flag("local-server-addr").Value.String(); localServerAddr != "" {
			objectConfig.LocalServerAddr = localServerAddr
		}
		if tunnelServerAddr := cmd.Flag("tunnel-server-addr").Value.String(); tunnelServerAddr != "" {
			objectConfig.TunnelServerAddr = tunnelServerAddr
		}
		if controllerAddr := cmd.Flag("controller-addr").Value.String(); controllerAddr != "" {
			objectConfig.ControllerAddr = controllerAddr
		}
		if publicServerAddr := cmd.Flag("public-server-addr").Value.String(); publicServerAddr != "" {
			objectConfig.PublicServerAddr = publicServerAddr
		}
		return nil
	}

	// 添加其他字段...
}
