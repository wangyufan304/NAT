package main

var object *objectConfigData

func initCobra() {
	object = new(objectConfigData)
	rootCmd.Flags().StringVarP(&object.Name, "name", "n", "", "Client name")
	rootCmd.Flags().StringVarP(&object.LocalServerAddr, "local-server-addr", "l", "", "The address of the local web server program")
	rootCmd.Flags().StringVarP(&object.TunnelServerAddr, "tunnel-server-addr", "t", "", "The address of the tunnel server used to connect the local and public networks")
	rootCmd.Flags().StringVarP(&object.ControllerAddr, "controller-addr", "c", "", "The address of the controller channel used to send controller messages to the client")
	rootCmd.Flags().StringVarP(&object.PublicServerAddr, "public-server-addr", "s", "", "The address of the public server used for accessing the inner web server")
	rootCmd.Flags().StringVarP(&object.UserName, "username", "u", "", "the name for auth the server.")
	rootCmd.Flags().StringVarP(&object.Password, "password", "P", "", "the password for auth the server.")
	// 添加其他字段...
}

func exchange() {
	if object.Name != "" {
		objectConfig.Name = object.Name
	}
	if object.TunnelServerAddr != "" {
		objectConfig.TunnelServerAddr = object.TunnelServerAddr
	}
	if object.ControllerAddr != "" {
		objectConfig.ControllerAddr = object.ControllerAddr
	}
	if object.PublicServerAddr != "" {
		objectConfig.PublicServerAddr = object.PublicServerAddr
	}
	if object.LocalServerAddr != "" {
		objectConfig.LocalServerAddr = object.LocalServerAddr
	}
	if object.UserName != "" {
		objectConfig.UserName = object.UserName
	}
	if object.Password != "" {
		objectConfig.Password = object.Password
	}
}
