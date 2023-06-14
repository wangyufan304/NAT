package main

import (
	"fmt"
	"github.com/byteYuFan/NAT/instance"
	"os"
	"strings"
)

func art() {
	fmt.Println("  _   _              _______    _____   ______   _____   __      __  ______   _____  \n | \\ | |     /\\     |__   __|  / ____| |  ____| |  __ \\  \\ \\    / / |  ____| |  __ \\ \n |  \\| |    /  \\       | |    | (___   | |__    | |__) |  \\ \\  / /  | |__    | |__) |\n | . ` |   / /\\ \\      | |     \\___ \\  |  __|   |  _  /    \\ \\/ /   |  __|   |  _  / \n | |\\  |  / ____ \\     | |     ____) | | |____  | | \\ \\     \\  /    | |____  | | \\ \\ \n |_| \\_| /_/    \\_\\    |_|    |_____/  |______| |_|  \\_\\     \\/     |______| |_|  \\_\\\n                                                                                     \n                                                                                     ")
	// fmt.Println("          _____                _____                _____          \n         /\\    \\              |\\    \\              /\\    \\         \n        /::\\____\\             |:\\____\\            /::\\    \\        \n       /::::|   |             |::|   |            \\:::\\    \\       \n      /:::::|   |             |::|   |             \\:::\\    \\      \n     /::::::|   |             |::|   |              \\:::\\    \\     \n    /:::/|::|   |             |::|   |               \\:::\\    \\    \n   /:::/ |::|   |             |::|   |               /::::\\    \\   \n  /:::/  |::|___|______       |::|___|______        /::::::\\    \\  \n /:::/   |::::::::\\    \\      /::::::::\\    \\      /:::/\\:::\\    \\ \n/:::/    |:::::::::\\____\\    /::::::::::\\____\\    /:::/  \\:::\\____\\\n\\::/    / ~~~~~/:::/    /   /:::/~~~~/~~         /:::/    \\::/    /\n \\/____/      /:::/    /   /:::/    /           /:::/    / \\/____/ \n             /:::/    /   /:::/    /           /:::/    /          \n            /:::/    /   /:::/    /           /:::/    /           \n           /:::/    /    \\::/    /            \\::/    /            \n          /:::/    /      \\/____/              \\/____/             \n         /:::/    /                                                \n        /:::/    /                                                 \n        \\::/    /                                                  \n         \\/____/                                                   \n                                                                   \n                                                                   \n                                                                   \n                                                                   \n                                                                   \n                                                                   \n                                                                   \n                                                                   \n                                                                   \n                                                                                                                         ")
}

func printServerRelationInformation() {
	// 打印信息
	fmt.Println("[ServerName]", objectConfig.Name)
	fmt.Println("[MaxServerConn]", objectConfig.MaxTCPConnNum)
	fmt.Println("[服务端开启端口]", objectConfig.ExposePort)
}

func initLogger() {
	flag := false
	if strings.ToLower(objectConfig.StartLog) == "true" {
		flag = true
	}
	myLogger = instance.NewMyLogger(os.Stdout, flag, objectConfig.LogFilename)
}
