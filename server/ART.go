package main

import "fmt"

func art() {
	fmt.Println("  _   _              _______    _____   ______   _____   __      __  ______   _____  \n | \\ | |     /\\     |__   __|  / ____| |  ____| |  __ \\  \\ \\    / / |  ____| |  __ \\ \n |  \\| |    /  \\       | |    | (___   | |__    | |__) |  \\ \\  / /  | |__    | |__) |\n | . ` |   / /\\ \\      | |     \\___ \\  |  __|   |  _  /    \\ \\/ /   |  __|   |  _  / \n | |\\  |  / ____ \\     | |     ____) | | |____  | | \\ \\     \\  /    | |____  | | \\ \\ \n |_| \\_| /_/    \\_\\    |_|    |_____/  |______| |_|  \\_\\     \\/     |______| |_|  \\_\\\n                                                                                     \n                                                                                     ")
}

func printServerRelationInformation() {
	// 打印信息
	fmt.Println("[CurrentConnInfo] ", serverInstance.CurrentConnInfo)
	fmt.Println("[Counter] ", serverInstance.Counter)
	fmt.Println("[TaskQueueSlice] ", serverInstance.TaskQueueSlice)
	fmt.Println("[MaxTCPConnSize] ", serverInstance.MaxTCPConnSize)
	fmt.Println("[MaxConnSize] ", serverInstance.MaxConnSize)
	fmt.Println("[ExposePort] ", serverInstance.ExposePort)
	fmt.Println("[TaskQueueBufferSize] ", serverInstance.TaskQueueBufferSize)
	fmt.Println("[TaskQueueSize] ", serverInstance.TaskQueueSize)
	fmt.Println("[ProcessingMap] ", serverInstance.ProcessingMap)
	fmt.Println("[ListenerAndClientConn] ", serverInstance.ListenerAndClientConn)

}
